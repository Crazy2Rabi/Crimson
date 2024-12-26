package agent

import (
	"context"
	"github.com/hsgames/gold/app"
	"github.com/jinzhu/now"
	"sync"
	"time"
)

var s *service

type service struct {
	sync.WaitGroup

	agents sync.Map

	ctx    context.Context
	cancel context.CancelFunc
}

func New() (app.Service, error) {
	s = &service{}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s, nil
}

func (s *service) Name() string {
	return "agent"
}

func (s *service) Init() error {
	return nil
}

func (s *service) Start() error {
	s.Add(1)
	defer s.Done()

	defer s.closeAgents()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case nowTime := <-ticker.C:
			s.checkAgents(nowTime)

			// todo 统计在线人数等等
		}
	}
}

func (s *service) Stop() error {
	s.cancel()
	s.Wait()
	return nil
}

func (s *service) closeAgents() {
	s.agents.Range(func(key, value any) bool {
		key.(*Agent).Conn.Close()
		return true
	})
}

func (s *service) checkAgents(nowTime time.Time) {
	s.agents.Range(func(key, value any) bool {
		a := key.(*Agent)
		openTime := a.openTime.Load().(now.Now).Time

		// 5秒内没有完成预登录，踢掉
		if !a.preLoginOk.Load() {
			if nowTime.Sub(openTime) >= 5*time.Second {
				a.Close(errPreLoginTimeout)
				s.agents.Delete(key)
				return true
			}
		}

		// 20秒内没登录进游戏，踢掉
		if !a.loginOk.Load() {
			if nowTime.Sub(openTime) >= 20*time.Second {
				a.Close(errLoginTimeout)
				s.agents.Delete(a)
				return true
			}
		}

		// 登录成功了，但是5分钟不进游戏，踢掉
		if !a.enterOk.Load() {
			if nowTime.Sub(openTime) >= 5*time.Minute {
				a.Close(errEnterTimeout)
				s.agents.Delete(a)
				return true
			}
		}

		return true
	})
}
