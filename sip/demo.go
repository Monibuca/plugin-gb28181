package sip

import "fmt"

func DemoMessage() {
	registerStr := `REGISTER sip:34020000002000000001@3402000000 SIP/2.0
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
	fmt.Println(registerStr)
	msg, err := Decode([]byte(registerStr))
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

func DemoVIA() {
	str1 := "SIP / 2.0 / UDP first.example.com: 4000;ttl=16 ;maddr=224.2.0.1 ;branch=z9hG4bKa7c6a8dlze.1"
	str2 := "SIP/2.0/UDP 192.168.1.64:5060;rport;received=192.168.1.64;branch=z9hG4bK1000615294"

	var err error
	v1 := &Via{}
	err = v1.Parse(str1)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	fmt.Printf("result:%v\n", v1.String())

	v2 := &Via{}
	err = v2.Parse(str2)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	fmt.Printf("result:%v\n", v2.String())

}
