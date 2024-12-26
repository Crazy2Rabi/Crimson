package def

import (
	"Common/Framework/dbredis"
	"github.com/gomodule/redigo/redis"
)

type FV struct {
	Field string
	Value any
}

type AccountInfo struct {
}

// 性别
type GenderType int8

const (
	GtNone GenderType = iota
	GtMale
	GtFemale
)

// 玩家
type Player struct {
	updateFVs []FV

	Uid        uint64
	Account    string // 账号
	Zone       uint64 // 区Id
	CreateTime int64  // 创建时间

	Name           string // 名字
	Gender         int8   // 性别
	ShowId         int32  // 显示Id
	LastLoginTime  int64  // 最后登录时间
	LastLoginIp    string // 最后登录IP
	LastLogoutTime int64  // 最后登出时间

	Gold    int32 // 金币
	Diamond int32 // 钻石
	Items   Items // 背包
}

func (p *Player) Load(key string) (err error) {
	conn := dbredis.Conn()
	defer conn.Close()

	var v []any

	if v, err = redis.Values(conn.Do("HGETALL", key)); err != nil {
		return
	}

	if len(v) == 0 {
		err = redis.ErrNil
		return
	}

	err = redis.ScanStruct(v, p)

	return
}

func (p *Player) Process(f func([]FV) error) error {
	defer func() {
		p.updateFVs = p.updateFVs[:0]
	}()

	for i, v := range p.updateFVs {
		switch x := v.Value.(type) {
		case redis.Argument:
			p.updateFVs[i].Value = x.RedisArg()
		}
	}

	return f(p.updateFVs)
}

func (p *Player) update(field string, value any) {
	for i := range p.updateFVs {
		if p.updateFVs[i].Field == field {
			p.updateFVs[i].Value = value
			return
		}
	}

	p.updateFVs = append(p.updateFVs, FV{field, value})
}

func (p *Player) UpdateUid(value uint64) {
	p.Uid = value
	p.update("Uid", value)
}

func (p *Player) UpdateAccount(account string) {
	p.Account = account
	p.update("Account", account)
}

func (p *Player) UpdateZone(zone uint64) {
	p.Zone = zone
	p.update("Zone", zone)
}

func (p *Player) UpdateCreateTime(value int64) {
	p.CreateTime = value
	p.update("CreateTime", value)
}

func (p *Player) UpdateLastLoginTime(value int64) {
	p.LastLoginTime = value
	p.update("LastLoginTime", value)
}

func (p *Player) UpdateLastLoginIp(value string) {
	p.LastLoginIp = value
	p.update("LastLoginIp", value)
}
