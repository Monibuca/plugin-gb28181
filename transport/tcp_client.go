package transport

import (
	"fmt"
	"net"
)

type TCPClient struct {
	Statistic
	host       string
	port       uint16
	conn       net.Conn
	readChan   chan *Packet
	writeChan  chan *Packet
	remoteAddr net.Addr
	localAddr  net.Addr
	done       chan struct{}
}

func NewTCPClient(host string, port uint16) IClient {
	return &TCPClient{
		host:      host,
		port:      port,
		readChan:  make(chan *Packet, 10),
		writeChan: make(chan *Packet, 10),
		done:      make(chan struct{}),
	}
}

func (c *TCPClient) IsReliable() bool {
	return true
}

func (c *TCPClient) Name() string {
	return fmt.Sprintf("tcp client to:%s", fmt.Sprintf("%s:%d", c.host, c.port))
}

func (c *TCPClient) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *TCPClient) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *TCPClient) Start() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		fmt.Println("dial tcp server failed :", err.Error())
		return err
	} else {
		fmt.Println("start tcp client")
	}

	c.conn = conn
	c.remoteAddr = conn.RemoteAddr()
	c.localAddr = conn.LocalAddr()

	//开启写线程
	go func() {
		for {
			select {
			case p := <-c.writeChan:
				_, err := c.conn.Write(p.Data)
				if err != nil {
					fmt.Println("client write failed:", err.Error())
					_ = c.Close()
					return
				}
			case <-c.done:
				return
			}
		}
	}()

	fmt.Println("start tcp client")
	fmt.Printf("remote addr: %s, local addr: %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())

	//读线程，阻塞
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("tcp client read error:", err.Error())
			return err
		}
		c.readChan <- &Packet{
			Addr: c.remoteAddr,
			Data: buf[:n],
		}
	}
}

func (c *TCPClient) ReadPacketChan() <-chan *Packet {
	return c.readChan
}

func (c *TCPClient) WritePacket(packet *Packet) {
	c.writeChan <- packet
}

func (c *TCPClient) Close() error {
	close(c.done)
	return c.conn.Close()
}

//外部定期调用此接口，实现心跳
func (c *TCPClient) Heartbeat(p *Packet) {
	if p == nil {
		p = &Packet{
			Data: []byte("ping"),
		}
	}
	c.WritePacket(p)
}
