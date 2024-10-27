package caddytraceid

import (
	"context"
	"net"

	"github.com/caddyserver/caddy/v2"
)

type Hello struct {
	Address string `json:"address"`

	listener net.Listener
}

func (Hello) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy_hello",
		New: func() caddy.Module { return new(Hello) },
	}
}

func init() {
	caddy.RegisterModule(Hello{})
}

func (h *Hello) Start() error {
	network, err := caddy.ParseNetworkAddress(h.Address)
	if err != nil {
		return err
	}

	listener, err := network.Listen(context.Background(), 0, net.ListenConfig{})
	if err != nil {
		return err
	}
	h.listener = listener.(net.Listener)
	go h.loop()
	return nil
}

func (h *Hello) Stop() error {
	return h.listener.Close()
}

func (h *Hello) loop() {
	for {
		c, err := h.listener.Accept()
		if err != nil {
			break
		}
		go h.handleConn(c)
	}
}

func (h *Hello) handleConn(c net.Conn) {
	_, _ = c.Write([]byte("Hello world"))
	_ = c.Close()
}
