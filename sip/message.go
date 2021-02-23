package sip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Monibuca/plugin-gb28181/v3/utils"
)

//Content-Type: Application/MANSCDP+xml
//Content-Type: Application/SDP
//Call-ID: 202081530679
//Max-Forwards: 70
//User-Agent: SRS/4.0.32(Leo)
//Subject: 34020000001320000001:0009093128,34020000002000000001:0
//Content-Length: 164

type Message struct {
	Mode Mode //0:REQUEST, 1:RESPONSE

	StartLine     *StartLine
	Via           *Via     //Via
	From          *Contact //From
	To            *Contact //To
	CallID        string   //Call-ID
	CSeq          *CSeq    //CSeq
	Contact       *Contact //Contact
	Authorization string   //Authorization
	MaxForwards   int      //Max-Forwards
	UserAgent     string   //User-Agent
	Subject       string   //Subject
	ContentType   string   //Content-Type
	Expires       int      //Expires
	ContentLength int      //Content-Length
	Route         *Contact
	Body          string
	Addr          string
}

func (m *Message) BuildResponse(code int) *Message {
	response := Message{
		Mode:        SIP_MESSAGE_RESPONSE,
		From:        m.From,
		To:          m.To,
		CallID:      m.CallID,
		CSeq:        m.CSeq,
		Via:         m.Via,
		MaxForwards: m.MaxForwards,
		StartLine: &StartLine{
			Code: code,
		},
	}
	return &response
}

//z9hG4bK + 10个随机数字
func randBranch() string {
	return fmt.Sprintf("z9hG4bK%s", utils.RandNumString(8))
}

func BuildMessageRequest(method Method, transport, sipSerial, sipRealm, username, srcIP string, srcPort uint16, expires, cseq int, body string) *Message {
	server := fmt.Sprintf("%s@%s", sipSerial, sipRealm)
	client := fmt.Sprintf("%s@%s", username, sipRealm)

	msg := &Message{
		Mode:          SIP_MESSAGE_REQUEST,
		MaxForwards:   70,
		UserAgent:     "IPC",
		Expires:       expires,
		ContentLength: 0,
	}
	msg.StartLine = &StartLine{
		Method: method,
		Uri:    NewURI(server),
	}
	msg.Via = &Via{
		Transport: transport,
		Host:      client,
	}
	msg.Via.Params = map[string]string{
		"branch": randBranch(),
		"rport":  "-1", //only key,no-value
	}
	msg.From = &Contact{
		Uri:    NewURI(client),
		Params: nil,
	}
	msg.From.Params = map[string]string{
		"tag": utils.RandNumString(10),
	}
	msg.To = &Contact{
		Uri: NewURI(client),
	}
	msg.CallID = utils.RandNumString(8)
	msg.CSeq = &CSeq{
		ID:     uint32(cseq),
		Method: method,
	}

	msg.Contact = &Contact{
		Uri: NewURI(fmt.Sprintf("%s@%s:%d", username, srcIP, srcPort)),
	}
	if len(body) > 0 {
		msg.ContentLength = len(body)
		msg.Body = body
	}
	return msg
}

func (m *Message) GetMode() Mode {
	return m.Mode
}

func (m *Message) IsRequest() bool {
	return m.Mode == SIP_MESSAGE_REQUEST
}

func (m *Message) IsResponse() bool {
	return m.Mode == SIP_MESSAGE_RESPONSE
}

func (m *Message) GetMethod() Method {
	if m.CSeq == nil {
		b, _ := Encode(m)
		println(string(b))
		return MESSAGE
	}
	return m.CSeq.Method
}

//此消息是否使用可靠传输
func (m *Message) IsReliable() bool {
	protocol := strings.ToUpper(m.Via.Transport)
	return "TCP" == protocol || "TLS" == protocol || "SCTP" == protocol
}

//response code
func (m *Message) GetStatusCode() int {
	return m.StartLine.Code
}

//response code and reason
func (m *Message) GetReason() string {
	return DumpError(m.StartLine.Code)
}

func (m *Message) GetBranch() string {
	if m.Via == nil {
		panic("invalid via")
	}
	if m.Via.Params == nil {
		panic("invalid via params")
	}

	b, ok := m.Via.Params["branch"]
	if !ok {
		panic("invalid via paramas branch")
	}

	return b
}

//构建响应消息的时候，会使用请求消息的 source 和 destination
//请求消息的source，格式： host:port
func (m *Message) Source() string {
	if m.Mode == SIP_MESSAGE_RESPONSE {
		fmt.Println("only for request message")
		return ""
	}

	if m.Via == nil {
		fmt.Println("invalid request message")
		return ""
	}

	var (
		host, port string
		via        = m.Via
	)

	if received, ok := via.Params["received"]; ok && received != "" {
		host = received
	} else {
		host = via.Host
	}

	if rport, ok := via.Params["rport"]; ok && rport != "-1" && rport != "0" && rport != "" {
		port = rport
	} else if via.Port != "" {
		port = via.Port
	} else {
		//如果port为空，则上层构建消息的时候，根据sip服务的默认端口来选择
	}
	return fmt.Sprintf("%v:%v", host, port)
}

//目标地址：这个应该是用于通过route头域实现proxy这样的功能，暂时不支持
func (m *Message) Destination() string {
	//TODO:
	return ""
}

//=======================================================================================================

func Decode(data []byte) (msg *Message, err error) {
	msg = &Message{}

	content := string(data)
	content = strings.Trim(content, CRLFCRLF)
	msgArr := strings.Split(content, CRLFCRLF)
	//第一部分：header
	//第二部分：body
	if len(msgArr) == 0 {
		fmt.Println("invalid sip message:", data)
		err = errors.New("invalid sip message")
		return
	}

	headStr := strings.TrimSpace(msgArr[0])
	if msgArrLen := len(msgArr); msgArrLen > 1 {
		for i := 1; i < msgArrLen; i++ {
			msg.Body += strings.TrimSpace(msgArr[i])
		}
	}

	headStr = strings.Trim(headStr, CRLF)
	headArr := strings.Split(headStr, CRLF)
	for i, line := range headArr {
		//fmt.Printf("%02d --- %s ---- %d\n", i, line, len(line))
		if i == 0 {
			firstline := strings.Trim(line, " ")
			tmp := strings.Split(firstline, " ")
			//if len(tmp) != 3 {
			//	fmt.Println("parse first line failed:", firstline)
			//	err = errors.New("invalid first line")
			//	return
			//}

			if strings.HasPrefix(firstline, VERSION) {
				//status line
				//SIP/2.0 200 OK
				var num int64
				num, err = strconv.ParseInt(tmp[1], 10, 64)
				if err != nil {
					return
				}
				msg.Mode = SIP_MESSAGE_RESPONSE
				msg.StartLine = &StartLine{
					raw:     firstline,
					Version: VERSION,
					Code:    int(num),
					phrase:  strings.Join(tmp[2:], " "),
				}
			} else {
				//request line
				//REGISTER sip:34020000002000000001@3402000000 SIP/2.0
				//MESSAGE sip:34020000002000000001@3402000000 SIP/2.0
				msg.Mode = SIP_MESSAGE_REQUEST
				msg.StartLine = &StartLine{
					raw:     firstline,
					Method:  Method(tmp[0]),
					Version: VERSION,
				}
				if len(tmp) > 1 {
					msg.StartLine.Uri, err = parseURI(tmp[1])
					if err != nil {
						return
					}
				} else {
					fmt.Println(firstline)
				}
			}
			continue
		}

		pos := strings.IndexByte(line, ':')
		if pos == -1 {
			continue
		}
		k := strings.ToLower(strings.TrimSpace(line[:pos]))
		v := strings.TrimSpace(line[pos+1:])
		//fmt.Printf("%02d ---k = %s , v = %s\n", i, k, v)

		if len(v) == 0 {
			continue
		}

		switch k {
		case "via":
			//Via: SIP/2.0/UDP 192.168.1.64:5060;rport;branch=z9hG4bK385701375
			msg.Via = &Via{}
			err = msg.Via.Parse(v)
			if err != nil {
				return
			}

		case "from":
			msg.From = &Contact{}
			err = msg.From.Parse(v)
			if err != nil {
				return
			}

		case "to":
			msg.To = &Contact{}
			err = msg.To.Parse(v)
			if err != nil {
				return
			}

		case "call-id":
			msg.CallID = v

		case "cseq":
			//CSeq: 2 REGISTER
			msg.CSeq = &CSeq{}
			err = msg.CSeq.Parse(v)
			if err != nil {
				return
			}

		case "contact":
			msg.Contact = &Contact{}
			err = msg.Contact.Parse(v)
			if err != nil {
				return
			}

		case "max-forwards":
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				fmt.Printf("parse head faield: %s,%s\n", k, v)
				return nil, err
			}
			msg.MaxForwards = int(n)

		case "user-agent":
			msg.UserAgent = v

		case "expires":
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				fmt.Printf("parse head faield: %s,%s\n", k, v)
				return nil, err
			}
			msg.Expires = int(n)

		case "content-length":
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				fmt.Printf("parse head faield: %s,%s\n", k, v)
				return nil, err
			}
			msg.ContentLength = int(n)

		case "authorization":
			msg.Authorization = v

		case "content-type":
			msg.ContentType = v
		case "route":
			//msg.Route = new(Contact)
			//msg.Route.Parse(v)
		default:
			fmt.Printf("invalid sip head: %s,%s\n", k, v)
		}
	}
	return
}

func Encode(msg *Message) ([]byte, error) {
	sb := strings.Builder{}
	sb.WriteString(msg.StartLine.String())
	sb.WriteString(CRLF)

	if msg.Via != nil {
		sb.WriteString("Via: ")
		sb.WriteString(msg.Via.String())
		sb.WriteString(CRLF)
	}

	if msg.From != nil {
		sb.WriteString("From: ")
		sb.WriteString(msg.From.String())
		sb.WriteString(CRLF)
	}

	if msg.To != nil {
		sb.WriteString("To: ")
		sb.WriteString(msg.To.String())
		sb.WriteString(CRLF)
	}

	if msg.CallID != "" {
		sb.WriteString("Call-ID: ")
		sb.WriteString(msg.CallID)
		sb.WriteString(CRLF)
	}
	if msg.CSeq != nil {
		sb.WriteString("CSeq: ")
		sb.WriteString(msg.CSeq.String())
		sb.WriteString(CRLF)
	}

	if msg.Contact != nil {
		sb.WriteString("Contact: ")
		sb.WriteString(msg.Contact.String())
		sb.WriteString(CRLF)
	}

	if msg.UserAgent != "" {
		sb.WriteString("User-Agent: ")
		sb.WriteString(msg.UserAgent)
		sb.WriteString(CRLF)
	}

	if msg.ContentType != "" {
		sb.WriteString("Content-Type: ")
		sb.WriteString(msg.ContentType)
		sb.WriteString(CRLF)
	}

	if msg.Expires != 0 {
		sb.WriteString("Expires: ")
		sb.WriteString(strconv.Itoa(msg.Expires))
		sb.WriteString(CRLF)
	}

	if msg.Subject != "" {
		sb.WriteString("Subject: ")
		sb.WriteString(msg.Subject)
		sb.WriteString(CRLF)
	}

	if msg.IsRequest() {
		//request only

		sb.WriteString("Max-Forwards: ")
		sb.WriteString(strconv.Itoa(msg.MaxForwards))
		sb.WriteString(CRLF)

		if msg.Authorization != "" {
			sb.WriteString("Authorization: ")
			sb.WriteString(msg.Authorization)
			sb.WriteString(CRLF)
		}
	} else {
		//response only
	}

	sb.WriteString("Content-Length: ")
	sb.WriteString(strconv.Itoa(msg.ContentLength))

	sb.WriteString(CRLFCRLF)

	if msg.Body != "" {
		sb.WriteString(msg.Body)
	}

	return []byte(sb.String()), nil
}
