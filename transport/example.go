package transport

import (
	"fmt"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	"os"
	"time"
)

//默认端口：TCP/UDP是5060，5061是在TCP上的TLS
//对于服务器监听UDP的任何端口和界面，都必须在TCP上也进行同样的监听。这是因为可能消息还需要通过TCP进行传输，比如消息过大的情况。
const SipHost string = "127.0.0.1"
const SipPort uint16 = 5060

func RunServerTCP() {
	tcp := NewTCPServer(SipPort, true)
	go PacketHandler(tcp)
	go func() {
		_ = tcp.Start()
	}()

	select {}
}

//测试通讯，客户端先发一条消息
func RunClientTCP() {
	c := NewTCPClient(SipHost, SipPort)
	go PacketHandler(c)
	go func() {
		_ = c.Start()
	}()

	//发送测试数据
	fmt.Println("send test data")
	go func() {
		for {
			c.WritePacket(&Packet{Data: []byte("from client : " + time.Now().String())})
			time.Sleep(2 * time.Second)
		}
	}()
	select {}
}
func PacketHandler(s ITransport) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("packet handler panic: ", err)
			utils.PrintStack()
			os.Exit(1)
		}
	}()

	fmt.Println("PacketHandler ========== ", s.Name())

	ch := s.ReadPacketChan()
	//阻塞读取消息
	for {
		select {
		case p := <-ch:
			fmt.Println("packet content:", string(p.Data))
			//TODO:message parse
		}
	}
}

//======================================================================

func RunServerUDP() {
	udp := NewUDPServer(SipPort)

	go PacketHandler(udp)
	go func() {
		_ = udp.Start()
	}()

	select {}
}

func RunClientUDP() {
	c := NewUDPClient(SipHost, SipPort)
	go PacketHandler(c)
	go func() {
		_ = c.Start()
	}()
	//发送测试数据
	go func() {
		for {
			time.Sleep(1 * time.Second)
			c.WritePacket(&Packet{
				Data: []byte("hello " + time.Now().String())})
		}
	}()

	select {}
}
