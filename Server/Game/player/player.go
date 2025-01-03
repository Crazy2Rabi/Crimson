package player

import (
	"Common/def"
	"net"
)

type Agent interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Push(e any)
	Shutdown(err error)
	Close(err error)
	LoginInfo() LoginInfo
	SetLoginInfo(LoginInfo)
}

type LoginInfo struct {
	Account string
	Token   string
}

func New(a Agent, account string, uid, zone uint64) *Player {
	return &Player{
		Agent: a,
		Player: &def.Player{
			Account: account,
			Uid:     uid,
			Zone:    zone,
		},
	}
}

type Player struct {
	Agent
	*def.Player

	EnterProcessing bool
	Temp            []any
}
