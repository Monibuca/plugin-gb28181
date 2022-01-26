package transaction

import (
	"context"
	"fmt"
	. "github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transport"
	"github.com/sirupsen/logrus"
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
	conn            *net.Conn
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

	c.tp.StartAndWait()
}

func (c *Core) handlerListen() {
	//阻塞读取消息
	for p := range c.tp.ReadPacketChan() {
		if len(p.Data) < 5 {
			continue
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
	// fmt.Println("packet content:", string(p.Data))
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

func (s *Core) NewTX(key string) *GBTx {
	return s.txs.NewTX(key, s.tp.UDPConn())
}
func (s *Core) GetTX(key string) *GBTx {
	return s.txs.GetTX(key)
}
func (s *Core) MustTX(key string) *GBTx {
	tx := s.txs.GetTX(key)
	if tx == nil {
		tx = s.NewTX(key)
	}
	return tx
}

func (s *Core) handlerRequest(msg *Request) {
	tx := s.MustTX(GetTXKey(msg.Message))
	logrus.Traceln("receive request from:", msg.Source(), ",method:", msg.GetMethod(), "txKey:", tx.Key(), "message: \n", msg.Event)
	s.hmu.RLock()
	handler, ok := s.requestHandlers[msg.GetMethod()]
	s.hmu.RUnlock()
	if !ok {
		logrus.Errorln("not found handler func,requestMethod:", msg.GetMethod(), msg.Event)
		go handlerMethodNotAllowed(msg, tx)
		return
	}

	go handler(msg, tx)
}

func (s *Core) handlerResponse(msg *Response) {
	tx := s.GetTX(GetTXKey(msg.Message))
	logrus.Traceln("receive response from:", msg.Source(), "txKey:", tx.Key(), "message: \n", msg.Event)
	if tx == nil {
		logrus.Infoln("not found tx. receive response from:", msg.Source(), "message: \n", msg.Event)
	} else {
		tx.ReceiveResponse(msg)
	}
}
func handlerMethodNotAllowed(req *Request, tx *GBTx) {
	var resp Response
	resp.Message = req.BuildResponse(http.StatusMethodNotAllowed)
	resp.DestAdd = req.SourceAdd
	resp.SourceAdd = req.DestAdd
	tx.Respond(&resp)
}

// Request Request
func (s *Core) Request(req *Request) (*GBTx, error) {
	var viaHop Via
	req.Via = &viaHop
	viaHop.Host = s.SipIP
	viaHop.Port = strconv.Itoa(int(s.SipPort))
	viaHop.Params["branch"] = RandBranch()
	viaHop.Params["rport"] = ""
	tx := s.MustTX(GetTXKey(req.Message))
	return tx, tx.Request(req)
}

// RegistHandler RegistHandler
func (s *Core) RegistHandler(method Method, handler RequestHandler) {
	s.hmu.Lock()
	s.requestHandlers[method] = handler
	s.hmu.Unlock()
}
