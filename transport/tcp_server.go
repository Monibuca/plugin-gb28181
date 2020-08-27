package transport

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	Statistic
	addr      string
	listener  net.Listener
	readChan  chan *Packet
	writeChan chan *Packet
	done      chan struct{}
	Keepalive bool
	sessions  sync.Map //key 是 remote-addr , value:*Connection。
}

func NewTCPServer(port uint16, keepalive bool) IServer {
	tcpAddr := fmt.Sprintf(":%d", port)

	return &TCPServer{
		addr:      tcpAddr,
		Keepalive: keepalive,
		readChan:  make(chan *Packet, 10),
		writeChan: make(chan *Packet, 10),
		done:      make(chan struct{}),
	}
}

func (s *TCPServer) IsReliable() bool {
	return true
}

func (s *TCPServer) Name() string {
	return fmt.Sprintf("tcp server at:%s", s.addr)
}
func (s *TCPServer) IsKeepalive() bool {
	return s.Keepalive
}

func (s *TCPServer) Start() error {
	//监听端口
	//开启tcp连接线程
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	//s.listener, err = tls.Listen("tcp", s.tcpAddr, tlsConfig)
	if err != nil {
		fmt.Println("TCP Listen failed:", err)
		return err
	}
	defer s.listener.Close()

	fmt.Println("start tcp server at: ", s.addr)

	//心跳线程
	if s.Keepalive {
		//TODO:start heartbeat thread
	}
	//写线程
	go func() {
		for {
			select {
			case p := <-s.writeChan:
				val, ok := s.sessions.Load(p.Addr.String())
				if !ok {
					return
				}
				c := val.(*Connection)
				_, _ = c.Conn.Write(p.Data)
			case <-s.done:
				return
			}

		}
	}()

	//读线程
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			var tempDelay time.Duration // how long to sleep on accept failure
			fmt.Println("accept err :", err.Error())
			//  重连。参考http server
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				time.Sleep(tempDelay)
				continue
			}
			fmt.Println("accept error, retry failed & exit.")
			return err
		}

		// conn.SetReadDeadline(time.Now().Add(600 * time.Second))
		session := &Connection{Conn: conn, Addr: conn.RemoteAddr()}
		address := session.Addr.String()
		s.sessions.Store(address, session)

		fmt.Println(fmt.Sprintf("new tcp client remoteAddr: %v", address))
		go s.handlerSession(session)
	}
}

func (s *TCPServer) handlerSession(c *Connection) {
	addrStr := c.Addr.String()

	//recovery from panic
	defer func() {
		s.CloseOne(addrStr)
		if err := recover(); err != nil {
			fmt.Println("client receiver handler panic: ", err)
		}
	}()

	buf := make([]byte, 2048)
	for {
		n, err := c.Conn.Read(buf)
		switch {
		case err == nil:
			p := &Packet{
				Addr: c.Addr,
				Data: buf[:n],
			}
			s.readChan <- p
		case err == io.EOF:
			fmt.Println(fmt.Sprintf("io.EOF,client close --- remoteAddr: %v", c.Addr))
			return
		case err != nil:
			fmt.Println("client other err: ", err)
			fmt.Println(fmt.Sprintf("client other err ---  remoteAddr: %v", addrStr))
			return
		}
	}
}

func (s *TCPServer) CloseOne(addr string) {
	val, ok := s.sessions.Load(addr)
	if !ok {
		return
	}
	c := val.(*Connection)
	_ = c.Conn.Close()
	s.sessions.Delete(addr)
}

func (s *TCPServer) ReadPacketChan() <-chan *Packet {
	return s.readChan
}
func (s *TCPServer) WritePacket(packet *Packet) {
	s.writeChan <- packet
}

func (s *TCPServer) Close() error {
	//TODO：TCP服务退出之前，需要先close掉所有客户端的连接
	s.sessions.Range(func(key, value interface{}) bool {
		c := value.(*Connection)
		_ = c.Conn.Close()
		s.sessions.Delete(key)
		return true
	})
	return nil
}
