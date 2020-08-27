package transport

import (
	"fmt"
	"net"
	"os"
)

type UDPClient struct {
	Statistic
	host       string
	port       uint16
	conn       *net.UDPConn
	readChan   chan *Packet
	writeChan  chan *Packet
	done       chan struct{}
	remoteAddr net.Addr
	localAddr  net.Addr
}

func NewUDPClient(host string, port uint16) IClient {
	return &UDPClient{
		host:      host,
		port:      port,
		readChan:  make(chan *Packet, 10),
		writeChan: make(chan *Packet, 10),
		done:      make(chan struct{}),
	}
}

func (c *UDPClient) IsReliable() bool {
	return false
}

func (c *UDPClient) Name() string {
	return fmt.Sprintf("udp client to:%s", fmt.Sprintf("%s:%d", c.host, c.port))
}
func (c *UDPClient) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *UDPClient) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *UDPClient) Start() error {
	addrStr := fmt.Sprintf("%s:%d", c.host, c.port)
	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	c.remoteAddr = conn.RemoteAddr()
	c.localAddr = conn.LocalAddr()

	fmt.Println("udp client remote addr:", conn.RemoteAddr().String())
	fmt.Println("udp client local addr:", conn.LocalAddr().String())

	//写线程
	go func() {
		for {
			select {
			case p := <-c.writeChan:
				_, err = conn.Write(p.Data)
				if err != nil {
					fmt.Println("udp client write failed:", err.Error())
					continue
				}
			case <-c.done:
				return
			}
		}
	}()

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err)
			os.Exit(1)
		}

		c.readChan <- &Packet{
			Addr: c.remoteAddr,
			Data: buf[:n],
		}
	}
}

func (c *UDPClient) ReadPacketChan() <-chan *Packet {
	return c.readChan
}

func (c *UDPClient) WritePacket(packet *Packet) {
	c.writeChan <- packet
}

func (c *UDPClient) Close() error {
	close(c.done)
	return c.conn.Close()
}

//外部定期调用此接口，实现心跳
func (c *UDPClient) Heartbeat(p *Packet) {
	if p == nil {
		p = &Packet{
			Data: []byte("ping"),
		}
	}
	c.WritePacket(p)
}
