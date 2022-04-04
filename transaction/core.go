package transaction

import (
	"context"
	"fmt"
	. "github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transport"
	. "github.com/Monibuca/utils/v3"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type RequestHandler func(req *Request, tx *GBTx)

//Core: transactions manager
//管理所有 transactions，以及相关全局参数、运行状态机
type Core struct {
	ctx             context.Context //上下文
	hmu             *sync.RWMutex
	requestHandlers map[Method]RequestHandler
	txs             *GBTxs
	tp              transport.ITransport //transport
	*Config                              //sip server配置信息
	OnRegister      func(msg *Request, tx *GBTx)
	OnMessage       func(msg *Request, tx *GBTx)
	udpaddr         net.Addr
}

//初始化一个 Core，需要能响应请求，也要能发起请求
//client 发起请求
//server 响应请求
//TODO:根据角色，增加相关配置信息
//TODO:通过context管理子线程
//TODO:单元测试
func NewCore(config *Config) *Core {
	ActiveTX = &GBTxs{
		Txs: map[string]*GBTx{},
		RWM: &sync.RWMutex{},
	}

	core := &Core{
		requestHandlers: map[Method]RequestHandler{},
		txs:             ActiveTX,
		Config:          config,
		ctx:             context.Background(),
		hmu:             &sync.RWMutex{},
	}

	if config.SipNetwork == "TCP" {
		core.tp = transport.NewTCPServer(config.SipPort, true)
	} else {
		core.tp = transport.NewUDPServer(config.SipPort)
	}

	return core
}

func (c *Core) StartAndWait() {
	go c.handlerListen()

	_ = c.tp.StartAndWait()
}

func (c *Core) handlerListen() {
	//阻塞读取消息
	for p := range c.tp.ReadPacketChan() {
		if len(p.Data) < 5 {
			continue
		}
		if c.LogVerbose {
			Println("Received: \n", string(p.Data))
		}
		if err := c.HandleReceiveMessage(p); err != nil {
			fmt.Println("handler sip response message failed:", err.Error())
			continue
		}
	}
}

//接收到的消息处理
//收到消息有两种：1、请求消息 2、响应消息
//请求消息则直接响应处理。
//响应消息则需要匹配到请求，让请求的transaction来处理。
//TODO：参考srs和osip的流程，以及文档，做最终处理。需要将逻辑分成两层：TU 层和 transaction 层
func (c *Core) HandleReceiveMessage(p *transport.Packet) (err error) {
	//Println("packet content:", string(p.Data))
	//var msg *Message
	msg, err := Decode(p.Data)
	if err != nil {
		fmt.Println("parse sip message failed:", err.Error())
		return ErrorParse
	}
	if msg.Via == nil || msg.From == nil {
		return ErrorParse
	}
	//这里不处理超过MTU的包，不处理半包
	err = checkMessage(msg)
	if err != nil {
		return err
	}

	if msg.IsRequest() {
		req := &Request{Message: msg}
		req.SourceAdd = p.Addr
		req.DestAdd = c.udpaddr
		c.handlerRequest(req)
	} else {
		//TODO:对于uac，收到response消息，是否要检查 rport 和 received 呢？因为uas可能对此做了修改
		resp := &Response{Message: msg}
		resp.SourceAdd = p.Addr
		resp.DestAdd = c.udpaddr
		c.handlerResponse(resp)
	}

	return
}

func (c *Core) NewTX(key string) *GBTx {
	tx := c.txs.NewTX(key, *c.tp.Conn())
	tx.Core = c
	return tx
}
func (c *Core) GetTX(key string) *GBTx {
	return c.txs.GetTX(key)
}
func (c *Core) MustTX(key string) *GBTx {
	tx := c.txs.GetTX(key)
	if tx == nil {
		tx = c.NewTX(key)
	}
	return tx
}

func (c *Core) handlerRequest(msg *Request) {
	tx := c.MustTX(GetTXKey(msg.Message))

	//Println("receive request from:", msg.Source(), ",method:", msg.GetMethod(), "txKey:", tx.Key(), "message: \n", string(encode))
	c.hmu.RLock()
	handler, ok := c.requestHandlers[msg.GetMethod()]
	c.hmu.RUnlock()
	if !ok {
		encode, _ := Encode(msg.Message)
		Println("not found handler func,requestMethod:", msg.GetMethod(), msg.Event, encode)
		go handlerMethodNotAllowed(msg, tx)
		return
	}

	go handler(msg, tx)
}

func (c *Core) handlerResponse(msg *Response) {
	tx := c.GetTX(GetTXKey(msg.Message))

	if tx == nil {
		str, _ := Encode(msg.Message)
		Println("not found tx. receive response from:", msg.Source(), "message: \n", string(str))
	} else {
		tx.ReceiveResponse(msg)
	}
}
func handlerMethodNotAllowed(req *Request, tx *GBTx) {
	var resp Response
	resp.Message = req.BuildResponse(http.StatusMethodNotAllowed)
	resp.DestAdd = req.SourceAdd
	resp.SourceAdd = req.DestAdd
	_ = tx.Respond(&resp)
}
func (c *Core) SipRequestForResponse(req *Request) (response *Response, err error) {
	var tx *GBTx
	tx, err = c.Request(req)
	if err == nil {
		return tx.SipResponse()
	}
	return
}
func (c *Core) ResolveAddress(msg *Message) (destAddr net.Addr, err error) {
	addr := msg.Addr

	if addr == "" {
		viaParams := msg.Via.Params
		var host, port string
		var ok1, ok2 bool
		if host, ok1 = viaParams["maddr"]; !ok1 {
			if host, ok2 = viaParams["received"]; !ok2 {
				host = msg.Via.Host
			}
		}
		//port
		port = viaParams["rport"]
		if port == "" || port == "0" || port == "-1" {
			port = msg.Via.Port
		}

		if port == "" {
			port = "5060"
		}

		addr = fmt.Sprintf("%s:%s", host, port)
	}

	// fmt.Println("dest addr:", addr)
	var err1, err2 error

	if msg.Via.Transport == "UDP" {
		destAddr, err2 = net.ResolveUDPAddr("udp", addr)
	} else {
		destAddr, err2 = net.ResolveTCPAddr("tcp", addr)
	}

	if err1 != nil {
		return nil, err1
	}

	if err2 != nil {
		return nil, err2
	}

	return destAddr, nil
}

// Request Request
func (c *Core) Request(req *Request) (*GBTx, error) {
	if req.Via == nil {
		var viaHop Via
		viaHop.Host = c.SipIP
		viaHop.Port = strconv.Itoa(int(c.SipPort))
		viaHop.Params = make(map[string]string)
		viaHop.Params["branch"] = RandBranch()
		viaHop.Params["rport"] = ""
		req.Via = &viaHop
	}
	tx := c.MustTX(GetTXKey(req.Message))
	return tx, tx.Request(req)
}

// Request Request
func (c *Core) Respond(resp *Response) (*GBTx, error) {

	tx := c.MustTX(GetTXKey(resp.Message))
	return tx, tx.Respond(resp)
}

// RegistHandler RegistHandler
func (c *Core) RegistHandler(method Method, handler RequestHandler) {
	c.hmu.Lock()
	c.requestHandlers[method] = handler
	c.hmu.Unlock()
}
