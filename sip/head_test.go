package sip

import (
	"fmt"
	"testing"
)

func TestContact(t *testing.T) {
	str1 := "\"Mr.Watson\" <sip:watson@worcester.bell-telephone.com>;q=0.7; expires=3600,\"Mr.Watson\" <mailto:watson@bell-telephone.com>"
	//str1 := `"Mr.Watson" <sip:watson@worcester.bell-telephone.com>;q=0.7;`
	c := &Contact{}
	err := c.Parse(str1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("source:", str1)
	fmt.Println("result:", c.String())
}

func TestVia(t *testing.T) {
	str1 := "SIP / 2.0 / UDP first.example.com: 4000;ttl=16 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1"
	str2 := "SIP/2.0/UDP 192.168.1.64:5060;rport;received=192.168.1.64;branch=z9hG4bK1000615294"

	var err error
	v1 := &Via{}
	err = v1.Parse(str1)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	fmt.Printf("source:%v\n", str1)
	fmt.Printf("result:%v\n", v1.String())

	v2 := &Via{}
	err = v2.Parse(str2)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	fmt.Printf("source:%v\n", str2)
	fmt.Printf("result:%v\n", v2.String())

}

func TestMessage1(t *testing.T) {
	str1 := `REGISTER sip:34020000002000000001@3402000000 SIP/2.0
Via: SIP/2.0/UDP 192.168.1.64:5060;rport;branch=z9hG4bK385701375
From: <sip:34020000001320000001@3402000000>;tag=1840661473
To: <sip:34020000001320000001@3402000000>
Call-ID: 418133739
CSeq: 1 REGISTER
Contact: <sip:34020000001320000001@192.168.1.64:5060>
Max-Forwards: 70
User-Agent: IP Camera
Expires: 3600
Content-Length: 0`

	fmt.Println("input:")
	fmt.Println(str1)
	msg, err := Decode([]byte(str1))
	if err != nil {
		fmt.Println("decode message failed:", err.Error())
		return
	}
	out, err := Encode(msg)
	if err != nil {
		fmt.Println("encode message failed:", err.Error())
		return
	}
	fmt.Println("=====================================")
	fmt.Println("output:")
	fmt.Println(string(out))
}

func TestMessage2(t *testing.T) {
	str1 := `SIP/2.0 200 OK
Via: SIP/2.0/UDP 192.168.1.151:5060;rport=5060;branch=SrsGbB56116414
From: <sip:34020000002000000001@3402000000>;tag=SrsGbF72006729
To: <sip:34020000001320000001@3402000000>;tag=416442565
Call-ID: 202093500940
CSeq: 101 INVITE
Contact: <sip:34020000001320000001@192.168.1.64:5060>
Content-Type: application/sdp
User-Agent: IP Camera
Content-Length:   185

v=0
o=34020000001320000001 1835 1835 IN IP4 192.168.1.64
s=Play
c=IN IP4 192.168.1.64
t=0 0
m=video 15060 RTP/AVP 96
a=sendonly
a=rtpmap:96 PS/90000
a=filesize:0
y=0009093131`

	fmt.Println("input:")
	fmt.Println(str1)
	msg, err := Decode([]byte(str1))
	if err != nil {
		fmt.Println("decode message failed:", err.Error())
		return
	}
	out, err := Encode(msg)
	if err != nil {
		fmt.Println("encode message failed:", err.Error())
		return
	}
	fmt.Println("=====================================")
	fmt.Println("output:")
	fmt.Println(string(out))
}
