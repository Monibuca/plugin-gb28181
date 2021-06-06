package transaction

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transport"
)

//状态机之状态
type State int

const (
	/* STATES for invite client transaction */
	ICT_PRE_CALLING State = iota
	ICT_CALLING
	ICT_PROCEEDING
	ICT_COMPLETED
	ICT_TERMINATED

	/* STATES for invite server transaction */
	IST_PRE_PROCEEDING
	IST_PROCEEDING
	IST_COMPLETED
	IST_CONFIRMED
	IST_TERMINATED

	/* STATES for NON-invite client transaction */
	NICT_PRE_TRYING
	NICT_TRYING
	NICT_PROCEEDING
	NICT_COMPLETED
	NICT_TERMINATED

	/* STATES for NON-invite server transaction */
	NIST_PRE_TRYING
	NIST_TRYING
	NIST_PROCEEDING
	NIST_COMPLETED
	NIST_TERMINATED

	/* STATES for dialog */
	DIALOG_EARLY
	DIALOG_CONFIRMED
	DIALOG_CLOSE
)

var stateMap = map[State]string{
	ICT_PRE_CALLING:    "ICT_PRE_CALLING",
	ICT_CALLING:        "ICT_CALLING",
	ICT_PROCEEDING:     "ICT_PROCEEDING",
	ICT_COMPLETED:      "ICT_COMPLETED",
	ICT_TERMINATED:     "ICT_TERMINATED",
	IST_PRE_PROCEEDING: "IST_PRE_PROCEEDING",
	IST_PROCEEDING:     "IST_PROCEEDING",
	IST_COMPLETED:      "IST_COMPLETED",
	IST_CONFIRMED:      "IST_CONFIRMED",
	IST_TERMINATED:     "IST_TERMINATED",
	NICT_PRE_TRYING:    "NICT_PRE_TRYING",
	NICT_TRYING:        "NICT_TRYING",
	NICT_PROCEEDING:    "NICT_PROCEEDING",
	NICT_COMPLETED:     "NICT_COMPLETED",
	NICT_TERMINATED:    "NICT_TERMINATED",
	NIST_PRE_TRYING:    "NIST_PRE_TRYING",
	NIST_TRYING:        "NIST_TRYING",
	NIST_PROCEEDING:    "NIST_PROCEEDING",
	NIST_COMPLETED:     "NIST_COMPLETED",
	NIST_TERMINATED:    "NIST_TERMINATED",
	DIALOG_EARLY:       "DIALOG_EARLY",
	DIALOG_CONFIRMED:   "DIALOG_CONFIRMED",
	DIALOG_CLOSE:       "DIALOG_CLOSE",
}

func (s State) String() string {
	return stateMap[s]
}

//状态机之事件
type Event int

const (
	/* TIMEOUT EVENTS for ICT */
	TIMEOUT_A Event = iota /**< Timer A */
	TIMEOUT_B              /**< Timer B */
	TIMEOUT_D              /**< Timer D */

	/* TIMEOUT EVENTS for NICT */
	TIMEOUT_E /**< Timer E */
	TIMEOUT_F /**< Timer F */
	TIMEOUT_K /**< Timer K */

	/* TIMEOUT EVENTS for IST */
	TIMEOUT_G /**< Timer G */
	TIMEOUT_H /**< Timer H */
	TIMEOUT_I /**< Timer I */

	/* TIMEOUT EVENTS for NIST */
	TIMEOUT_J /**< Timer J */

	/* FOR INCOMING MESSAGE */
	RCV_REQINVITE     /**< Event is an incoming INVITE request */
	RCV_REQACK        /**< Event is an incoming ACK request */
	RCV_REQUEST       /**< Event is an incoming NON-INVITE and NON-ACK request */
	RCV_STATUS_1XX    /**< Event is an incoming informational response */
	RCV_STATUS_2XX    /**< Event is an incoming 2XX response */
	RCV_STATUS_3456XX /**< Event is an incoming final response (not 2XX) */

	/* FOR OUTGOING MESSAGE */
	SND_REQINVITE     /**< Event is an outgoing INVITE request */
	SND_REQACK        /**< Event is an outgoing ACK request */
	SND_REQUEST       /**< Event is an outgoing NON-INVITE and NON-ACK request */
	SND_STATUS_1XX    /**< Event is an outgoing informational response */
	SND_STATUS_2XX    /**< Event is an outgoing 2XX response */
	SND_STATUS_3456XX /**< Event is an outgoing final response (not 2XX) */

	KILL_TRANSACTION /**< Event to 'kill' the transaction before termination */
	UNKNOWN_EVT      /**< Max event */
)

var eventMap = map[Event]string{
	TIMEOUT_A:         "TIMEOUT_A",
	TIMEOUT_B:         "TIMEOUT_B",
	TIMEOUT_D:         "TIMEOUT_D",
	TIMEOUT_E:         "TIMEOUT_E",
	TIMEOUT_F:         "TIMEOUT_F",
	TIMEOUT_K:         "TIMEOUT_K",
	TIMEOUT_G:         "TIMEOUT_G",
	TIMEOUT_H:         "TIMEOUT_H",
	TIMEOUT_I:         "TIMEOUT_I",
	TIMEOUT_J:         "TIMEOUT_J",
	RCV_REQINVITE:     "RCV_REQINVITE",
	RCV_REQACK:        "RCV_REQACK",
	RCV_REQUEST:       "RCV_REQUEST",
	RCV_STATUS_1XX:    "RCV_STATUS_1XX",
	RCV_STATUS_2XX:    "RCV_STATUS_2XX",
	RCV_STATUS_3456XX: "RCV_STATUS_3456XX",
	SND_REQINVITE:     "SND_REQINVITE",
	SND_REQACK:        "SND_REQACK",
	SND_REQUEST:       "SND_REQUEST",
	SND_STATUS_1XX:    "SND_STATUS_1XX",
	SND_STATUS_2XX:    "SND_STATUS_2XX",
	SND_STATUS_3456XX: "SND_STATUS_3456XX",
	KILL_TRANSACTION:  "KILL_TRANSACTION",
	UNKNOWN_EVT:       "UNKNOWN_EVT",
}

func (e Event) String() string {
	return eventMap[e]
}

//incoming SIP MESSAGE
func (e Event) IsIncomingMessage() bool {
	return e >= RCV_REQINVITE && e <= RCV_STATUS_3456XX
}

//incoming SIP REQUEST
func (e Event) IsIncomingRequest() bool {
	return e == RCV_REQINVITE || e == RCV_REQACK || e == RCV_REQUEST
}

//incoming SIP RESPONSE
func (e Event) IsIncomingResponse() bool {
	return e == RCV_STATUS_1XX || e == RCV_STATUS_2XX || e == RCV_STATUS_3456XX
}

//outgoing SIP MESSAGE
func (e Event) IsOutgoingMessage() bool {
	return e >= SND_REQINVITE && e <= SND_REQINVITE
}

//outgoing SIP REQUEST
func (e Event) IsOutgoingRequest() bool {
	return e == SND_REQINVITE || e == SND_REQACK || e == SND_REQUEST
}

//outgoing SIP RESPONSE
func (e Event) IsOutgoingResponse() bool {
	return e == SND_STATUS_1XX || e == SND_STATUS_2XX || e == SND_STATUS_3456XX
}

//a SIP MESSAGE
func (e Event) IsSipMessage() bool {
	return e >= RCV_REQINVITE && e <= SND_STATUS_3456XX
}

type EventObj struct {
	evt Event  // event type
	tid string // transaction id
	msg *sip.Message
}

//状态机类型
type FSMType int

const (
	FSM_ICT     FSMType = iota /**< Invite Client (outgoing) Transaction */
	FSM_IST                    /**< Invite Server (incoming) Transaction */
	FSM_NICT                   /**< Non-Invite Client (outgoing) Transaction */
	FSM_NIST                   /**< Non-Invite Server (incoming) Transaction */
	FSM_UNKNOWN                /**< Invalid Transaction */
)

var typeMap = map[FSMType]string{
	FSM_ICT:     "FSM_ICT",
	FSM_IST:     "FSM_IST",
	FSM_NICT:    "FSM_NICT",
	FSM_NIST:    "FSM_NIST",
	FSM_UNKNOWN: "FSM_UNKNOWN",
}

func (t FSMType) String() string {
	return typeMap[t]
}

//对外将sip通讯封装成请求和响应
//TODO：可参考http的request和response，屏蔽sip协议细节
type Request struct {
	data *sip.Message
}

//Code = 0，则响应正常
//Code != 0，打印错误提示信息 Message
type Response struct {
	Code    int
	Message string
	Data    *sip.Message
}

type Handler func(t *Transaction, e *EventObj) error //操作

type Header map[string]string

// timer相关基础常量、方法等定义
const (
	T1      = 100 * time.Millisecond
	T2      = 4 * time.Second
	T4      = 5 * time.Second
	TimeA   = T1
	TimeB   = 64 * T1
	TimeD   = 32 * time.Second
	TimeE   = T1
	TimeF   = 64 * T1
	TimeG   = T1
	TimeH   = 64 * T1
	TimeI   = T4
	TimeJ   = 64 * T1
	TimeK   = T4
	Time1xx = 100 * time.Millisecond
)

//TODO：是否要管理当前 transaction 的多次请求和响应的message？
//TODO：是否要管理当前 transaction 的头域
//TODO：多种transaction在一个struct里面管理不太方便，暂时写在一起，后期重构分开，并使用interface 解耦

//是否需要tp layer？
type Transaction struct {
	ctx        context.Context //线程管理、其他参数
	id         string          //transaction ID
	isReliable bool            //是否可靠传输
	core       *Core           //全局参数
	typo       FSMType         //状态机类型
	done       chan struct{}   //主动退出

	state    State          //当前状态
	event    chan *EventObj //输入的事件，带缓冲
	response chan *Response //输出的响应
	startAt  time.Time      //开始时间
	endAt    time.Time      //结束时间

	//messages []*sip.Message  //传输的消息缓存，origin request/last response/request ack...
	//header       Header //创建事物的消息头域参数:Via From  To CallID CSeq
	via          *sip.Via
	from         *sip.Contact
	to           *sip.Contact
	callID       string
	cseq         *sip.CSeq
	origRequest  *sip.Message //Initial request
	lastResponse *sip.Message //Last response，可能是临时的，也可能是最终的
	ack          *sip.Message //ack request sent

	//timer for ict
	timerA *SipTimer
	timerB *time.Timer
	timerD *time.Timer

	//timer for nict
	timerE *SipTimer
	timerF *time.Timer
	timerK *time.Timer

	//timer for ist
	timerG *time.Timer
	timerH *time.Timer
	timerI *time.Timer

	//timer for nist
	timerJ *time.Timer
}

type SipTimer struct {
	tm      *time.Timer
	timeout time.Duration //当前超时时间
	max     time.Duration //最大超时时间
}

func NewSipTimer(d, max time.Duration, f func()) *SipTimer {
	return &SipTimer{
		tm:      time.AfterFunc(d, f),
		timeout: d,
		max:     max,
	}
}

func (t *SipTimer) Reset(d time.Duration) {
	t.timeout = d
	if t.timeout > t.max && t.max != 0 {
		t.timeout = t.max
	}
	t.tm.Reset(t.timeout)
}

func (ta *Transaction) SetState(s State) {
	ta.state = s
}

func (ta *Transaction) GetTid() string {
	return ta.id
}

//每一个transaction至少有一个状态机线程运行
//TODO:如果是一个uac的transaction，则把最后响应的消息返回（通过response chan）
//transaction有很多消息需要传递到TU，也接收来自TU的消息。
func (ta *Transaction) Run() {
	for {
		select {
		case e := <-ta.event:
			//根据event调用对应的handler
			//fmt.Println("fsm run event:", e.evt.String())
			core := ta.core
			state := ta.state
			evtHandlers, ok1 := core.handlers[state]
			if !ok1 {
				//fmt.Println("invalid state:", ta.state.String())
				break
			}
			f, ok2 := evtHandlers[e.evt]
			if !ok2 {
				//fmt.Println("invalid handler for this event:", e.evt.String())
				break
			}
			//fmt.Printf("state:%s, event:%s\n", state.String(), e.evt.String())
			err := f(ta, e)
			if err != nil {
				//fmt.Printf("transaction run failed, state:%s, event:%s\n", state.String(), e.evt.String())
			}
		case <-ta.done:
			//fmt.Println("fsm exit")
			return

		case <-ta.ctx.Done():
			//fmt.Println("fsm killed")
			return
		}
	}
}

//Terminated:事物的终止
//TODO：check调用时机
func (ta *Transaction) Terminate() {

	ta.state = NICT_TERMINATED

	switch ta.typo {
	case FSM_ICT:
		ta.state = ICT_TERMINATED
	case FSM_NICT:
		ta.state = NICT_TERMINATED
	case FSM_IST:
		ta.state = IST_TERMINATED
	case FSM_NIST:
		ta.state = NIST_TERMINATED
	}

	//关掉事物的线程
	close(ta.done)

	//TODO：某些timer需要检查并关掉，并且设置为nil

	//remove ta from core
	ta.core.removeTa <- ta.id
}

//根据sip消息，解析出目标服务器地址，发送消息
func (ta *Transaction) SipSend(msg *sip.Message) error {
	err := checkMessage(msg)
	if err != nil {
		return err
	}
	addr := msg.Addr
	if addr == "" {
		viaParams := msg.Via.Params
		//host
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

	//fmt.Println("dest addr:", addr)

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
	ta.core.tp.WritePacket(pkt)

	return nil
}
