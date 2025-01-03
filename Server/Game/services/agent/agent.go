package agent

import (
	"Common/Framework/codec"
	"Common/message"
	"Game/manager"
	"Game/player"
	"errors"
	"fmt"
	"github.com/hsgames/gold/net"
	"github.com/jinzhu/now"
	"log/slog"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var (
	errUnsupportedEvent = errors.New("agent: unsupported event")
	errNotPreLogin      = errors.New("agent: not pre login")
	errNotLogin         = errors.New("agent: not login")
	errNotEnter         = errors.New("agent: not enter")
	errAlreadyLogin     = errors.New("agent: already login")
	errPreLoginTimeout  = errors.New("agent: pre login timeout")
	errLoginTimeout     = errors.New("agent: login timeout")
	errEnterTimeout     = errors.New("agent: enter timeout")
	errChanIsFull       = errors.New("agent: chan is full")
	errPlayerIsNil      = errors.New("agent: player is full")
)

type readEvent struct {
	data []byte
}

type Agent struct {
	net.Conn

	ch   chan any
	done chan struct{}

	preLoginOk atomic.Bool
	loginOk    atomic.Bool
	enterOk    atomic.Bool
	openTime   atomic.Value

	loginInfo player.LoginInfo
	player    *player.Player
}

func NewAgent(conn net.Conn) *Agent {
	a := &Agent{
		Conn: conn,
		ch:   make(chan any, 64),
		done: make(chan struct{}),
	}

	a.openTime.Store(*now.New(time.Now()))

	return a
}

func (a *Agent) OnOpen() error {
	s.agents.Store(a, struct{}{})

	// todo service 管理协程
	s.Add(1)
	go func() {
		defer s.Done()

		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("agent: run panic [%v]", r)
				a.Close(err)
				debug.PrintStack()
			}
		}()

		if err := a.run(); err != nil {
			err = fmt.Errorf("agent: run error [%w]", err)
			a.Close(err)
			debug.PrintStack()
		}
	}()

	return nil
}

func (a *Agent) OnClose() {
	a.Push(nil)

	s.agents.Delete(a)
}

func (a *Agent) OnRead(data []byte) error {
	a.Push(readEvent{data: data})
	return nil
}

func (a *Agent) Push(e any) {
	select {
	case a.ch <- e:
	default:
		close(a.done)
		a.Close(errChanIsFull)
	}
}

func (a *Agent) Write(m interface{}) error {
	data, err := codec.Encode(m)
	if err != nil {
		return err
	}

	return a.Conn.Write(data)
}

func (a *Agent) Shutdown(err error) {
	slog.Info("agent: shutdown",
		slog.String("agent", a.String()),
		slog.Any("error", err))
	a.Conn.Shutdown()
}

func (a *Agent) Close(err error) {
	slog.Info("agent: close",
		slog.String("agent", a.String()),
		slog.Any("error", err))
	a.Conn.Close()
}

func (a *Agent) LoginInfo() player.LoginInfo {
	return a.loginInfo
}

func (a *Agent) SetLoginInfo(loginInfo player.LoginInfo) {
	a.loginInfo = loginInfo
}

func (a *Agent) run() (err error) {
	if err = a.onOpen(); err != nil {
		return
	}

	defer a.onClose()

	for {
		select {
		case <-a.done:
			return
		case e := <-a.ch:
			if e == nil {
				return
			}

			if err = a.onEvent(e); err != nil {
				return
			}
		}
	}
}

func (a *Agent) onEvent(e any) (err error) {
	switch x := e.(type) {
	case readEvent:
		var m interface{}
		if m, err = codec.Decode(x.data); err != nil {
			return
		}

		err = a.onMessage(m)
	default:
		err = errUnsupportedEvent
	}
	return
}

func (a *Agent) onOpen() (err error) {
	slog.Debug("agent: open", slog.String("agent", a.String()))
	return
}

func (a *Agent) onClose() {
	slog.Debug("agent: close", slog.String("agent", a.String()))

	if a.player != nil {
		// todo 下线、存储
	}
}

func (a *Agent) onMessage(m interface{}) (err error) {
	if !a.preLoginOk.Load() {
		switch req := m.(type) {
		case message.PreLoginReq:
			slog.Debug("agent: OnPreLoginReq", slog.Any("req", req))

			var res message.PreLoginRes

			if res, err = manager.OnPreLoginReq(a, &req); err != nil {
				err = fmt.Errorf("agent: OnPreLoginReq error [%w]", err)
				return
			}

			if err = a.Write(res); err != nil {
				err = fmt.Errorf("agent: write error [%w]", err)
				return
			}

			slog.Debug("agent: OnPreLoginReq", slog.Any("res", res))

			a.preLoginOk.Store(true)
		default:
			err = errNotPreLogin
		}
		return
	}

	if !a.loginOk.Load() {
		switch req := m.(type) {
		case message.LoginReq:
			slog.Debug("agent: OnLoginReq", slog.Any("req", req))

			var res message.LoginRes

			if res, err = manager.OnLoginReq(a, &req); err != nil {
				err = fmt.Errorf("agent: OnLoginReq error [%w]", err)
				return
			}

			// todo 登录逻辑

			if err = a.Write(res); err != nil {
				err = fmt.Errorf("agent: OnLoginReq error [%w]", err)
			}

			a.loginOk.Store(true)
		default:
			err = errNotLogin
		}
		return
	}

	if !a.enterOk.Load() {
		switch req := m.(type) {
		case message.EnterReq:
			slog.Debug("agent: OnEnterReq", slog.Any("req", req))

			var (
				res message.EnterRes
				p   *player.Player
			)

			if p, res, err = manager.OnEnterReq(a, &req); err != nil {
				err = fmt.Errorf("agent: OnEnterReq error [%w]", err)
				return
			}

			if err = a.Write(res); err != nil {
				err = fmt.Errorf("agent: write error [%w]", err)
				return
			}

			a.enterOk.Store(true)
			a.player = p
		default:
			err = errNotEnter
		}
		return
	}

	if a.player == nil {
		err = errPlayerIsNil
		return
	}

	switch /*req :=*/ m.(type) {
	case message.PreLoginReq, message.LoginReq, message.EnterReq:
		err = errAlreadyLogin
	default:
		// todo 一般消息处理
	}

	return
}
