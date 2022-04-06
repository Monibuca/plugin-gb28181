package gb28181

import (
	"context"
	"net"

	"m7s.live/engine/v4/config"
	"m7s.live/engine/v4/log"
)

type UDP struct {
	ListenAddr string
	ListenNum  int //同时并行监听数量，0为CPU核心数量
}
type UDPPlugin interface {
	config.Plugin
	ServeUDP(*net.UDPConn)
}

func (udp *UDP) Listen(ctx context.Context, plugin UDPPlugin) error {
	addr, _ := net.ResolveUDPAddr("udp", udp.ListenAddr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("%s: Listen error: %v", udp.ListenAddr, err)
		return err
	}
	go plugin.ServeUDP(conn)
	<-ctx.Done()
	return nil
}
