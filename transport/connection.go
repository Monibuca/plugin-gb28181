package transport

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Connection Wrapper around net.Conn.
type Connection interface {
	net.Conn
	Network() string
	// String() string
	ReadFrom(buf []byte) (num int, raddr net.Addr, err error)
	WriteTo(buf []byte, raddr net.Addr) (num int, err error)
}

// Connection implementation.
type connection struct {
	baseConn       net.Conn
	laddr          net.Addr
	raddr          net.Addr
	mu             sync.RWMutex
	logKey         string
	Online         bool
	ReconnectCount int64 //重连次数
}

func newUDPConnection(baseConn net.Conn) Connection {
	conn := &connection{
		baseConn: baseConn,
		laddr:    baseConn.LocalAddr(),
		raddr:    baseConn.RemoteAddr(),
		logKey:   "udpConnection",
	}
	return conn
}
func newTCPConnection(baseConn net.Conn) Connection {
	conn := &connection{
		baseConn: baseConn,
		laddr:    baseConn.LocalAddr(),
		raddr:    baseConn.RemoteAddr(),
		logKey:   "udpConnection",
	}
	return conn
}

func (conn *connection) Read(buf []byte) (int, error) {
	var (
		num int
		err error
	)

	num, err = conn.baseConn.Read(buf)

	return num, err
}

func (conn *connection) ReadFrom(buf []byte) (num int, raddr net.Addr, err error) {
	num, raddr, err = conn.baseConn.(net.PacketConn).ReadFrom(buf)
	if err != nil {
		return num, raddr, err
	}
	fmt.Printf("readFrom %d , %s -> %s \n %s", num, raddr, conn.LocalAddr(), string(buf[:num]))
	return num, raddr, err
}

func (conn *connection) Write(buf []byte) (int, error) {
	var (
		num int
		err error
	)
	num, err = conn.baseConn.Write(buf)
	return num, err
}

func (conn *connection) WriteTo(buf []byte, raddr net.Addr) (num int, err error) {
	num, err = conn.baseConn.(net.PacketConn).WriteTo(buf, raddr)
	if err != nil {
		return num, err
	}
	//Printf("writeTo %d , %s -> %s \n %s", num, conn.baseConn.LocalAddr(), raddr.String(), string(buf[:num]))
	return num, err
}

func (conn *connection) LocalAddr() net.Addr {
	return conn.baseConn.LocalAddr()
}

func (conn *connection) RemoteAddr() net.Addr {
	return conn.baseConn.RemoteAddr()
}

func (conn *connection) Close() error {
	err := conn.baseConn.Close()
	return err
}

func (conn *connection) Network() string {
	return strings.ToUpper(conn.baseConn.LocalAddr().Network())
}

func (conn *connection) SetDeadline(t time.Time) error {
	return conn.baseConn.SetDeadline(t)
}

func (conn *connection) SetReadDeadline(t time.Time) error {
	return conn.baseConn.SetReadDeadline(t)
}

func (conn *connection) SetWriteDeadline(t time.Time) error {
	return conn.baseConn.SetWriteDeadline(t)
}
