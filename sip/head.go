package sip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//换行符号：
//linux,unix :     \r\n
//windows    :     \n
//Mac OS     :     \r
const (
	VERSION  = "SIP/2.0"  // sip version
	CRLF     = "\r\n"     // 0x0D0A
	CRLFCRLF = "\r\n\r\n" // 0x0D0A0D0A

	//CRLF     = "\n"   // 0x0D
	//CRLFCRLF = "\n\n" // 0x0D0D
)

//SIP消息类型：请求or响应
type Mode int

const (
	SIP_MESSAGE_REQUEST  Mode = 0
	SIP_MESSAGE_RESPONSE Mode = 1
)

//sip request method
type Method string

const (
	ACK       Method = "ACK"
	BYE       Method = "BYE"
	CANCEL    Method = "CANCEL"
	INVITE    Method = "INVITE"
	OPTIONS   Method = "OPTIONS"
	REGISTER  Method = "REGISTER"
	NOTIFY    Method = "NOTIFY"
	SUBSCRIBE Method = "SUBSCRIBE"
	MESSAGE   Method = "MESSAGE"
	REFER     Method = "REFER"
	INFO      Method = "INFO"
	PRACK     Method = "PRACK"
	UPDATE    Method = "UPDATE"
	PUBLISH   Method = "PUBLISH"
)

//startline
//MESSAGE sip:34020000001320000001@3402000000 SIP/2.0
//SIP/2.0 200 OK
type StartLine struct {
	raw string //原始内容

	//request line: method uri version
	Method  Method
	Uri     URI //Request-URI:请求的服务地址，不能包含空白字符或者控制字符，并且禁止用”<>”括上。
	Version string

	//status line: version code phrase
	Code   int //status code
	phrase string
}

func (l *StartLine) String() string {
	if l.Version == "" {
		l.Version = "SIP/2.0"
	}
	var result string
	if l.Method == "" {
		result = fmt.Sprintf("%s %d %s", l.Version, l.Code, l.phrase)
	} else {
		result = fmt.Sprintf("%s %s %s", l.Method, l.Uri.String(), l.Version)
	}
	l.raw = result
	return l.raw
}

//To From Referto Contact
//From: <sip:34020000001320000001@3402000000>;tag=575945878
//To: <sip:34020000002000000001@3402000000>
//Contact: <sip:34020000001320000001@27.38.49.149:49243>
//Contact: <sip:34020000001320000001@192.168.1.64:5060>;expires=0
type Contact struct {
	raw string //原始内容

	Nickname string //可以没有
	Uri      URI    //

	//header params
	Params map[string]string // include tag/q/expires
}

func (c *Contact) String() string {
	sb := strings.Builder{}

	if c.Nickname != "" {
		sb.WriteByte('"')
		sb.WriteString(c.Nickname)
		sb.WriteByte('"')
		sb.WriteByte(' ')
	}
	urlStr := c.Uri.String()
	if strings.ContainsAny(urlStr, ",?:") {
		urlStr = fmt.Sprintf("<%s>", urlStr)
	}
	sb.WriteString(urlStr)

	if c.Params != nil {
		for k, v := range c.Params {
			sb.WriteString(";")
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
		}
	}

	c.raw = sb.String()
	return c.raw
}

func (c *Contact) Parse(str string) (err error) {
	c.raw = str

	if str == "*" {
		c.Uri.host = "*"
		return
	}

	n0 := strings.IndexByte(str, '"')
	if n0 != -1 {
		str = str[n0+1:]
		n1 := strings.IndexByte(str, '"')
		if n1 == -1 {
			return errors.New("parse nickname failed")
		}
		c.Nickname = str[:n1]
		str = strings.TrimSpace(str[n1+1:])
	}

	if len(str) == 0 {
		return
	}

	var uriDone = false
	if strings.ContainsAny(str, "<>") {
		n2 := strings.IndexByte(str, '<')
		n3 := strings.IndexByte(str, '>')
		if n2 == -1 || n3 == -1 {
			err = errors.New("parse contact-uri failed")
			return
		}
		c.Uri, err = parseURI(str[n2+1 : n3])
		if err != nil {
			return
		}
		uriDone = true
		str = strings.TrimSpace(str[n3+1:])
	}

	if len(str) == 0 {
		return
	}

	str = strings.Trim(str, ";")
	arr1 := strings.Split(str, ";")
	for idx, one := range arr1 {
		//如果上面没有通过<>解析出来uri，则用分号split的第一个元素，就是uri字符串
		if !uriDone && idx == 0 {
			c.Uri, err = parseURI(one)
			if err != nil {
				return
			}

			continue
		}
		if c.Params == nil {
			c.Params = make(map[string]string)
		}
		arr2 := strings.Split(one, "=")
		k, v := arr2[0], arr2[1]
		c.Params[k] = v
	}

	return
}

//Via: SIP/2.0/UDP 192.168.1.64:5060;rport=49243;received=27.38.49.149;branch=z9hG4bK879576192
//Params：
//Received : IPv4address / IPv6address
//RPort    : 0-not found, -1-no-value, other-value
//Branch   : branch参数的值必须用magic cookie "z9hG4bK" 作为开头

/*
Via               =  ( "Via" / "v" ) HCOLON via-parm *(COMMA via-parm)
via-parm          =  sent-protocol LWS sent-by *( SEMI via-params )
via-params        =  via-ttl / via-maddr
                     / via-received / via-branch
                     / via-extension
via-ttl           =  "ttl" EQUAL ttl
via-maddr         =  "maddr" EQUAL host
via-received      =  "received" EQUAL (IPv4address / IPv6address)
via-branch        =  "branch" EQUAL token
via-extension     =  generic-param
sent-protocol     =  protocol-name SLASH protocol-version
                     SLASH transport
protocol-name     =  "SIP" / token
protocol-version  =  token
transport         =  "UDP" / "TCP" / "TLS" / "SCTP"
                     / other-transport
sent-by           =  host [ COLON port ]
ttl               =  1*3DIGIT ; 0 to 255
*/
type Via struct {
	raw       string // 原始内容
	Version   string // sip version: default to  SIP/2.0
	Transport string // UDP,TCP ,TLS , SCTP
	Host      string // sent-by : host:port
	Port      string //
	//header params
	Params map[string]string // include branch/rport/received/ttl/maddr
}

func (v *Via) GetBranch() string {
	return v.Params["branch"]
}

func (v *Via) GetSendBy() string {
	var host, port string

	sb := strings.Builder{}
	received := v.Params["received"]
	rport := v.Params["rport"]

	if received != "" {
		host = received
	} else {
		host = v.Host
	}

	if rport != "" && rport != "0" && rport != "-1" {
		port = rport
	} else if v.Port != "" {
		port = v.Port
	} else {
		if strings.ToUpper(v.Transport) == "UDP" {
			port = "5060"
		} else {
			port = "5061"
		}
	}

	sb.WriteString(host)
	sb.WriteString(":")
	sb.WriteString(port)

	return sb.String()
}
func (v *Via) String() string {
	sb := strings.Builder{}
	if v.Version == "" {
		v.Version = "SIP/2.0"
	}

	if v.Transport == "" {
		v.Transport = "UDP"
	}

	sb.WriteString(v.Version)
	sb.WriteString("/")
	sb.WriteString(v.Transport)
	sb.WriteString(" ")
	sb.WriteString(v.Host)
	if v.Port != "" {
		sb.WriteString(":")
		sb.WriteString(v.Port)
	}

	if v.Params != nil {
		for k, v := range v.Params {
			sb.WriteString(";")
			sb.WriteString(k)
			if v == "-1" {
				//rport 值为-1的时候，没有值
				continue
			}
			sb.WriteString("=")
			sb.WriteString(v)
		}
	}

	v.raw = sb.String()
	return v.raw
}

//注意via允许以下这种添加空白
//Via: SIP / 2.0 / UDP first.example.com: 4000;ttl=16 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1
//Via: SIP/2.0/UDP 192.168.1.64:5060;rport=5060;received=192.168.1.64;branch=z9hG4bK1000615294
func (v *Via) Parse(str string) (err error) {
	v.raw = str

	str = strings.Trim(str, ";")
	arr1 := strings.Split(str, ";")
	part1 := strings.TrimSpace(arr1[0]) //SIP / 2.0 / UDP first.example.com: 4000

	v.Host, v.Port = "", ""
	if n1 := strings.IndexByte(part1, ':'); n1 != -1 {
		v.Port = strings.TrimSpace(part1[n1+1:])
		part1 = strings.TrimSpace(part1[:n1])
	}

	n2 := strings.LastIndexByte(part1, ' ')
	if n2 == -1 {
		v.Host = part1 //error?
	} else {
		v.Host = strings.TrimSpace(part1[n2+1:])

		//解析protocol、version和transport，SIP / 2.0 / UDP
		part2 := part1[:n2]
		arr2 := strings.Split(part2, "/")
		if len(arr2) != 3 {
			err = errors.New("parse contait part1.1 failed:" + part2)
			return
		}
		v.Version = fmt.Sprintf("%s/%s", strings.TrimSpace(arr2[0]), strings.TrimSpace(arr2[1]))
		v.Transport = strings.TrimSpace(arr2[2])
	}

	//必须有参数
	v.Params = make(map[string]string)
	for i, one := range arr1 {
		if i == 0 {
			//arr[0]已经处理
			continue
		}
		one = strings.TrimSpace(one)
		arr2 := strings.Split(one, "=")
		//rport 这个参数可能没有 value。 -1:no-value, other-value
		if len(arr2) == 1 {
			if arr2[0] == "rport" {
				v.Params["rport"] = "-1"
				continue
			} else {
				fmt.Println("invalid param:", one)
				continue
			}
		}

		k, val := arr2[0], arr2[1]
		v.Params[k] = val
	}

	return
}

//CSeq: 101 INVITE
//CSeq: 2 REGISTER
type CSeq struct {
	raw    string //原始内容
	ID     uint32
	Method Method
}

func (c *CSeq) String() string {
	c.raw = fmt.Sprintf("%d %s", c.ID, c.Method)
	return c.raw
}

func (c *CSeq) Parse(str string) error {
	c.raw = str
	arr1 := strings.Split(str, " ")
	n, err := strconv.ParseInt(arr1[0], 10, 64)
	if err != nil {
		fmt.Println("parse cseq faield:", str)
		return err
	}
	c.ID = uint32(n)
	c.Method = Method(arr1[1])
	return nil
}

//sip:user:password@domain;uri-parameters?headers
/*
RFC3261
SIP-URI          =  "sip:" [ userinfo ] hostport
                    uri-parameters [ headers ]
SIPS-URI         =  "sips:" [ userinfo ] hostport
                    uri-parameters [ headers ]
userinfo         =  ( user / telephone-subscriber ) [ ":" password ] "@"
user             =  1*( unreserved / escaped / user-unreserved )
user-unreserved  =  "&" / "=" / "+" / "$" / "," / ";" / "?" / "/"
password         =  *( unreserved / escaped /
                    "&" / "=" / "+" / "$" / "," )
hostport         =  host [ ":" port ]
host             =  hostname / IPv4address / IPv6reference
hostname         =  *( domainlabel "." ) toplabel [ "." ]
domainlabel      =  alphanum
                    / alphanum *( alphanum / "-" ) alphanum
toplabel         =  ALPHA / ALPHA *( alphanum / "-" ) alphanum

IPv4address    =  1*3DIGIT "." 1*3DIGIT "." 1*3DIGIT "." 1*3DIGIT
IPv6reference  =  "[" IPv6address "]"
IPv6address    =  hexpart [ ":" IPv4address ]
hexpart        =  hexseq / hexseq "::" [ hexseq ] / "::" [ hexseq ]
hexseq         =  hex4 *( ":" hex4)
hex4           =  1*4HEXDIG
port           =  1*DIGIT

uri-parameters    =  *( ";" uri-parameter)
uri-parameter     =  transport-param / user-param / method-param
                     / ttl-param / maddr-param / lr-param / other-param
transport-param   =  "transport="
                     ( "udp" / "tcp" / "sctp" / "tls"
                     / other-transport)
other-transport   =  token
user-param        =  "user=" ( "phone" / "ip" / other-user)
other-user        =  token
method-param      =  "method=" Method
ttl-param         =  "ttl=" ttl
maddr-param       =  "maddr=" host
lr-param          =  "lr"
other-param       =  pname [ "=" pvalue ]
pname             =  1*paramchar
pvalue            =  1*paramchar
paramchar         =  param-unreserved / unreserved / escaped
param-unreserved  =  "[" / "]" / "/" / ":" / "&" / "+" / "$"

headers         =  "?" header *( "&" header )
header          =  hname "=" hvalue
hname           =  1*( hnv-unreserved / unreserved / escaped )
hvalue          =  *( hnv-unreserved / unreserved / escaped )
hnv-unreserved  =  "[" / "]" / "/" / "?" / ":" / "+" / "$"
*/
type URI struct {
	scheme  string            // sip sips
	host    string            // userinfo@domain  or userinfo@ip:port
	method  string            // uri和method有关？
	params  map[string]string // include branch/maddr/received/ttl/rport
	headers map[string]string // include branch/maddr/received/ttl/rport
}

func (u *URI) Host() string {
	return u.host
}
func (u *URI) UserInfo() string {
	return strings.Split(u.host, "@")[0]
}
func (u *URI) Domain() string {
	return strings.Split(u.host, "@")[1]
}
func (u *URI) IP() string {
	t := strings.Split(u.host, "@")
	if len(t) == 1 {
		return strings.Split(t[0], ":")[0]
	}
	return strings.Split(t[1], ":")[0]
}
func (u *URI) Port() string {
	t := strings.Split(u.host, "@")
	if len(t) == 1 {
		return strings.Split(t[0], ":")[1]
	}
	return strings.Split(t[1], ":")[1]
}
func (u *URI) String() string {
	if u.scheme == "" {
		u.scheme = "sip"
	}
	sb := strings.Builder{}
	sb.WriteString(u.scheme)
	sb.WriteString(":")
	sb.WriteString(u.host)
	if u.params != nil {
		for k, v := range u.params {
			sb.WriteString(";")
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
		}
	}

	if u.headers != nil {
		sb.WriteString("?")
		for k, v := range u.headers {
			sb.WriteString("&")
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
		}
	}

	return sb.String()
}

//对于gb28181，request-uri 不带参数
func NewURI(host string) URI {
	return URI{
		scheme: "sip",
		host:   host,
	}
}
func parseURI(str string) (ret URI, err error) {
	ret = URI{}

	//解析scheme
	str = strings.TrimSpace(str)
	n1 := strings.IndexByte(str, ':')
	if n1 == -1 {
		err = errors.New("invalid sheme")
		return
	}
	ret.scheme = str[:n1]
	str = str[n1+1:]
	if len(str) == 0 {
		return
	}

	//解析host
	n2 := strings.IndexByte(str, ';')
	if n2 == -1 {
		ret.host = str
		return
	}
	ret.host = str[:n2]

	str = str[n2+1:]
	if len(str) == 0 {
		return
	}

	//解析params and headers
	var paramStr, headerStr = "", ""
	n3 := strings.IndexByte(str, '?')
	if n3 == -1 {
		paramStr = str
	} else {
		paramStr = str[:n3]
		headerStr = str[n3+1:]
	}

	//k1=v1;k2=v2
	if paramStr != "" {
		ret.params = make(map[string]string)
		paramStr = strings.Trim(paramStr, ";")
		arr1 := strings.Split(paramStr, ";")
		for _, one := range arr1 {
			tmp := strings.Split(one, "=")
			if len(tmp) == 2 {
				k, v := tmp[0], tmp[1]
				ret.params[k] = v
			} else {
				ret.params[tmp[0]] = ""
			}
		}
	}

	//k1=v1&k2=v2
	if headerStr != "" {
		ret.headers = make(map[string]string)
		arr2 := strings.Split(paramStr, "&")
		for _, one := range arr2 {
			tmp := strings.Split(one, "=")
			k, v := tmp[0], tmp[1]
			ret.headers[k] = v
		}
	}

	return
}
