package server

import (
	"Game/services/agent"
	"github.com/hsgames/gold/net"
)

type handler struct {
	service *service
	agent   *agent.Agent
}

func (h *handler) OnOpen(conn net.Conn) error {
	h.agent = agent.NewAgent(conn)
	return h.agent.OnOpen()
}

func (h *handler) OnClose(conn net.Conn) {
	h.agent.OnClose()
}

func (h *handler) OnRead(conn net.Conn, data []byte) error {
	return h.agent.OnRead(data)
}
