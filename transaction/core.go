package transaction

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transport"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
)

//Core: transactions manager
//管理所有 transactions，以及相关全局参数、运行状态机
type Core struct {
	ctx          context.Context             //上下文
	handlers     map[State]map[Event]Handler //每个状态都可以处理有限个事件。不必加锁。
	transactions map[string]*Transaction     //管理所有 transactions,key:tid,value:transaction
	mutex        sync.RWMutex                //transactions的锁
	removeTa     chan string                 //要删除transaction的时候，通过chan传递tid
	tp           transport.ITransport        //transport
	*Config                                  //sip server配置信息
	OnRegister   func(*sip.Message)
	OnMessage    func(*sip.Message) bool
}

//初始化一个 Core，需要能响应请求，也要能发起请求
//client 发起请求
//server 响应请求
//TODO:根据角色，增加相关配置信息
//TODO:通过context管理子线程
//TODO:单元测试
func NewCore(config *Config) *Core {
	core := &Core{
		handlers:     make(map[State]map[Event]Handler),
		transactions: make(map[string]*Transaction),
		removeTa:     make(chan string, 10),
		Config:       config,
		ctx:          context.Background(),
	}
	if config.SipNetwork == "TCP" {
		core.tp = transport.NewTCPServer(config.SipPort, true)
	} else {
		core.tp = transport.NewUDPServer(config.SipPort)
	}
	//填充fsm
	core.addICTHandler()
	core.addISTHandler()
	core.addNICTHandler()
	core.addNISTHandler()
	return core
}

//add transaction to core
func (c *Core) AddTransaction(ta *Transaction) {
	c.mutex.Lock()
	c.transactions[ta.id] = ta
	c.mutex.Unlock()
	go ta.Run()
}

//delete transaction
func (c *Core) DelTransaction(tid string) {
	c.mutex.Lock()
	delete(c.transactions, tid)
	c.mutex.Unlock()
}

//创建事件:根据接收到的消息创建消息事件
func (c *Core) NewInComingMessageEvent(m *sip.Message) *EventObj {
	return &EventObj{
		evt: getInComingMessageEvent(m),
		tid: getMessageTransactionID(m),
		msg: m,
	}
}

//创建事件:根据发出的消息创建消息事件
func (c *Core) NewOutGoingMessageEvent(m *sip.Message) *EventObj {
	return &EventObj{
		evt: getOutGoingMessageEvent(m),
		tid: getMessageTransactionID(m),
		msg: m,
	}
}

//创建事物
//填充此事物的参数：via、from、to、callID、cseq
func (c *Core) initTransaction(ctx context.Context, obj *EventObj) *Transaction {
	m := obj.msg

	//ack要么属于一个invite事物，要么由TU层直接管理，不通过事物管理。
	if m.GetMethod() == sip.ACK {
		fmt.Println("ack nerver create transaction")
		return nil
	}
	ta := &Transaction{
		id:       obj.tid,
		core:     c,
		ctx:      ctx,
		done:     make(chan struct{}),
		event:    make(chan *EventObj, 10), //带缓冲的event channel
		response: make(chan *Response),
		startAt:  time.Now(),
		endAt:    time.Now().Add(1000000 * time.Hour),
	}
	//填充其他transaction的信息
	ta.via = m.Via
	ta.from = m.From
	ta.to = m.To
	ta.callID = m.CallID
	ta.cseq = m.CSeq
	ta.origRequest = m

	return ta
}

//状态机初始化:ICT
func (c *Core) addICTHandler() {
	c.addHandler(ICT_PRE_CALLING, SND_REQINVITE, ict_snd_invite)
	c.addHandler(ICT_CALLING, TIMEOUT_A, osip_ict_timeout_a_event)
	c.addHandler(ICT_CALLING, TIMEOUT_B, osip_ict_timeout_b_event)
	c.addHandler(ICT_CALLING, RCV_STATUS_1XX, ict_rcv_1xx)
	c.addHandler(ICT_CALLING, RCV_STATUS_2XX, ict_rcv_2xx)
	c.addHandler(ICT_CALLING, RCV_STATUS_3456XX, ict_rcv_3456xx)
	c.addHandler(ICT_PROCEEDING, RCV_STATUS_1XX, ict_rcv_1xx)
	c.addHandler(ICT_PROCEEDING, RCV_STATUS_2XX, ict_rcv_2xx)
	c.addHandler(ICT_PROCEEDING, RCV_STATUS_3456XX, ict_rcv_3456xx)
	c.addHandler(ICT_COMPLETED, RCV_STATUS_3456XX, ict_retransmit_ack)
	c.addHandler(ICT_COMPLETED, TIMEOUT_D, osip_ict_timeout_d_event)
}

//状态机初始化:IST
func (c *Core) addISTHandler() {
	c.addHandler(IST_PRE_PROCEEDING, RCV_REQINVITE, ist_rcv_invite)
	c.addHandler(IST_PROCEEDING, RCV_REQINVITE, ist_rcv_invite)
	c.addHandler(IST_COMPLETED, RCV_REQINVITE, ist_rcv_invite)
	c.addHandler(IST_COMPLETED, TIMEOUT_G, osip_ist_timeout_g_event)
	c.addHandler(IST_COMPLETED, TIMEOUT_H, osip_ist_timeout_h_event)
	c.addHandler(IST_PROCEEDING, SND_STATUS_1XX, ist_snd_1xx)
	c.addHandler(IST_PROCEEDING, SND_STATUS_2XX, ist_snd_2xx)
	c.addHandler(IST_PROCEEDING, SND_STATUS_3456XX, ist_snd_3456xx)
	c.addHandler(IST_COMPLETED, RCV_REQACK, ist_rcv_ack)
	c.addHandler(IST_CONFIRMED, RCV_REQACK, ist_rcv_ack)
	c.addHandler(IST_CONFIRMED, TIMEOUT_I, osip_ist_timeout_i_event)
}

//状态机初始化:NICT
func (c *Core) addNICTHandler() {
	c.addHandler(NICT_PRE_TRYING, SND_REQUEST, nict_snd_request)
	c.addHandler(NICT_TRYING, TIMEOUT_F, osip_nict_timeout_f_event)
	c.addHandler(NICT_TRYING, TIMEOUT_E, osip_nict_timeout_e_event)
	c.addHandler(NICT_TRYING, RCV_STATUS_1XX, nict_rcv_1xx)
	c.addHandler(NICT_TRYING, RCV_STATUS_2XX, nict_rcv_23456xx)
	c.addHandler(NICT_TRYING, RCV_STATUS_3456XX, nict_rcv_23456xx)
	c.addHandler(NICT_PROCEEDING, TIMEOUT_F, osip_nict_timeout_f_event)
	c.addHandler(NICT_PROCEEDING, TIMEOUT_E, osip_nict_timeout_e_event)
	c.addHandler(NICT_PROCEEDING, RCV_STATUS_1XX, nict_rcv_1xx)
	c.addHandler(NICT_PROCEEDING, RCV_STATUS_2XX, nict_rcv_23456xx)
	c.addHandler(NICT_PROCEEDING, RCV_STATUS_3456XX, nict_rcv_23456xx)
	c.addHandler(NICT_COMPLETED, TIMEOUT_K, osip_nict_timeout_k_event)
}

//状态机初始化:NIST
func (c *Core) addNISTHandler() {
	c.addHandler(NIST_PRE_TRYING, RCV_REQUEST, nist_rcv_request)
	c.addHandler(NIST_TRYING, SND_STATUS_1XX, nist_snd_1xx)
	c.addHandler(NIST_TRYING, SND_STATUS_2XX, nist_snd_23456xx)
	c.addHandler(NIST_TRYING, SND_STATUS_3456XX, nist_snd_23456xx)
	c.addHandler(NIST_PROCEEDING, SND_STATUS_1XX, nist_snd_1xx)
	c.addHandler(NIST_PROCEEDING, SND_STATUS_2XX, nist_snd_23456xx)
	c.addHandler(NIST_PROCEEDING, SND_STATUS_3456XX, nist_snd_23456xx)
	c.addHandler(NIST_PROCEEDING, RCV_REQUEST, nist_rcv_request)
	c.addHandler(NIST_COMPLETED, TIMEOUT_J, osip_nist_timeout_j_event)
	c.addHandler(NIST_COMPLETED, RCV_REQUEST, nist_rcv_request)
}

//状态机初始化：根据state 匹配到对应的状态机
func (c *Core) addHandler(state State, event Event, handler Handler) {
	m := c.handlers

	if state >= DIALOG_CLOSE {
		fmt.Println("invalid state:", state)
		return
	}

	if event >= UNKNOWN_EVT {
		fmt.Println("invalid event:", event)
		return
	}

	if _, ok := m[state]; !ok {
		m[state] = make(map[Event]Handler)
	}

	if _, ok := m[state][event]; ok {
		fmt.Printf("state:%d,event:%d, has been exist\n", state, event)
	} else {
		m[state][event] = handler
	}
}

func (c *Core) Start() {
	go c.Handler()

	c.tp.Start()
}

func (c *Core) Handler() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("packet handler panic: ", err)
			utils.PrintStack()
			os.Exit(1)
		}
	}()
	ch := c.tp.ReadPacketChan()
	//阻塞读取消息
	for {
		//fmt.Println("PacketHandler ========== SIP Client")
		select {
		case tid := <-c.removeTa:
			c.DelTransaction(tid)
		case p := <-ch:
			err := c.HandleReceiveMessage(p)
			if err != nil {
				fmt.Println("handler sip response message failed:", err.Error())
				continue
			}
		}
	}
}

//发送消息：发送请求或者响应
//发送消息仅负责发送。报错有两种：1、发送错误。2、发送了但是超时没有收到响应
//如果发送成功，如何判断是否收到响应？没有收到响应要重传
//所以一个transaction 有read和wriet的chan。
//发送的时候写 write chan
//接收的时候读取 read chan
//发送之后，就开启timer，超时重传，还要记录和修改每次超时时间。不超时的话，记得删掉timer
//发送 register 消息
func (c *Core) SendMessage(msg *sip.Message) *Response {
	method := msg.GetMethod()
	// data, _ := sip.Encode(msg)
	// fmt.Println("send message:", method)

	e := c.NewOutGoingMessageEvent(msg)

	//匹配事物
	c.mutex.RLock()
	ta, ok := c.transactions[e.tid]
	c.mutex.RUnlock()
	if !ok {
		//新的请求
		ta = c.initTransaction(c.ctx, e)

		//如果是sip 消息事件，则将消息缓存，填充typo和state
		if msg.IsRequest() {
			//as uac
			if method == sip.INVITE || method == sip.ACK {
				ta.typo = FSM_ICT
				ta.state = ICT_PRE_CALLING
			} else {
				ta.typo = FSM_NICT
				ta.state = NICT_PRE_TRYING
			}
		} else {
			//as uas:send response

		}

		c.AddTransaction(ta)
	}

	//把event推到transaction
	ta.event <- e
	<-ta.done
	if ta.lastResponse != nil {
		return &Response{
			Code:    ta.lastResponse.GetStatusCode(),
			Data:    ta.lastResponse,
			Message: ta.lastResponse.GetReason(),
		}
	} else {
		return &Response{
			Code: 504,
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
	var msg *sip.Message
	msg, err = sip.Decode(p.Data)
	if err != nil {
		fmt.Println("parse sip message failed:", err.Error())
		return ErrorParse
	}
	if msg.Via == nil {
		return ErrorParse
	}
	//这里不处理超过MTU的包，不处理半包
	err = checkMessage(msg)
	if err != nil {
		return err
	}

	//fmt.Println("receive message:", msg.GetMethod())

	e := c.NewInComingMessageEvent(msg)

	//一般应该是uas对于接收到的request做预处理
	if msg.IsRequest() {
		fixReceiveMessageViaParams(msg, p.Addr)
	} else {
		//TODO:对于uac，收到response消息，是否要检查 rport 和 received 呢？因为uas可能对此做了修改
	}
	//TODO：CANCEL、BYE 和 ACK 需要特殊处理，使用事物或者直接由TU层处理
	//查找transaction
	c.mutex.RLock()
	ta, ok := c.transactions[e.tid]
	c.mutex.RUnlock()
	method := msg.GetMethod()
	if msg.IsRequest() {
		switch method {
		case sip.ACK:
			//TODO:this should be a ACK for 2xx (but could be a late ACK!)
			return
		case sip.BYE:
			c.Send(msg.BuildResponse(200))
			return
		case sip.MESSAGE:
			if c.OnMessage(msg) && ta == nil {
				c.Send(msg.BuildResponse(200))
			}
			if ta != nil {
				ta.event <- c.NewOutGoingMessageEvent(msg.BuildResponse(200))
			}
		case sip.REGISTER:
			if !ok {
				ta = c.initTransaction(c.ctx, e)
				ta.typo = FSM_NIST
				ta.state = NIST_PROCEEDING
				c.AddTransaction(ta)
			}
			c.OnRegister(msg)
			ta.event <- c.NewOutGoingMessageEvent(msg.BuildResponse(200))
		//case sip.INVITE:
		//	ta.typo = FSM_IST
		//	ta.state = IST_PRE_PROCEEDING
		case sip.CANCEL:
			//TODO:CANCEL处理
			/* special handling for CANCEL */
			/* in the new spec, if the CANCEL has a Via branch, then it
			is the same as the one in the original INVITE */
			return
		}
	} else if ok {
		ta.event <- e

	}
	//TODO：TU层处理：根据需要，创建，或者匹配 Dialog
	//通过tag匹配到call和dialog
	//处理是否要重传ack
	return
}
func (c *Core) Send(msg *sip.Message) error {
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
	pkt := &transport.Packet{}
	pkt.Data, err1 = sip.Encode(msg)

	if msg.Via.Transport == "UDP" {
		pkt.Addr, err2 = net.ResolveUDPAddr("udp", addr)
	} else {
		pkt.Addr, err2 = net.ResolveTCPAddr("tcp", addr)
	}

	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}
	c.tp.WritePacket(pkt)
	return nil
}
