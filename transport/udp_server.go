package transport

import (
	"fmt"
	"net"
	"os"
)

type UDPServer struct {
	Statistic
	addr      string
	conn      *net.UDPConn
	readChan  chan *Packet
	writeChan chan *Packet
	done      chan struct{}
	Keepalive bool
	//Sessions  sync.Map // key is remote-addr的string , value:*Connection。UDP不需要
}

func NewUDPServer(port uint16) IServer {
	addrStr := fmt.Sprintf(":%d", port)

	return &UDPServer{
		addr:      addrStr,
		readChan:  make(chan *Packet, 1024),
		writeChan: make(chan *Packet, 1024),
		done:      make(chan struct{}),
	}
}

func (s *UDPServer) IsReliable() bool {
	return false
}

func (s *UDPServer) Name() string {
	return fmt.Sprintf("udp client to:%s", s.addr)
}

func (s *UDPServer) IsKeepalive() bool {
	return s.Keepalive
}

func (s *UDPServer) StartAndWait() error {
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listenUDP :", err)
		os.Exit(1)
	}
	defer func() {
		_ = conn.Close()
	}()
	ccc := *conn

	s.conn = &ccc

	fmt.Println("start udp server at: ", s.addr)

	//心跳线程
	if s.Keepalive {
		//TODO:start heartbeat thread
	}

	//读线程
	for {
		data := make([]byte, 4096)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err.Error())
			continue
		}
		s.readChan <- &Packet{
			Addr: remoteAddr,
			Data: data[:n],
		}
	}
}
func (s *UDPServer) ReadPacketChan() <-chan *Packet {
	return s.readChan
}
func (s *UDPServer) WritePacket(packet *Packet) {
	s.writeChan <- packet
}

func (s *UDPServer) Close() error {
	//所有session离线和关闭处理
	return nil
}
func (s *UDPServer) CloseOne(addr string) {
	//处理某设备离线
}
func (s *UDPServer) UDPConn() *net.UDPConn {
	return s.conn
}
func (s *UDPServer) Conn() *net.Conn {
	return nil
}
