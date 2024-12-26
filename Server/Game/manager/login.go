package manager

import (
	"Common/Framework/config"
	"Common/Framework/dbredis"
	"Common/Framework/uid"
	"Common/message"
	"Game/player"
	"errors"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
	"time"
)

const preLoginCheckCode uint32 = 0x32388545

var (
	errCheckCodeMismatch = errors.New("manager: check code is mismatch")
	errTokenMismatch     = errors.New("manager: token is mismatch")
	errTokenExpired      = errors.New("manager: token is expired")
)

func OnPreLoginReq(a player.Agent, req *message.PreLoginReq) (
	res message.PreLoginRes, err error) {
	if req.CheckCode != preLoginCheckCode {
		err = errCheckCodeMismatch
		return
	}

	res = message.PreLoginRes{}

	return
}

func getAccountKey(account string) string {
	return fmt.Sprintf("%sAccount:{%d}:{%s}",
		config.Instance().App.Prefix, getAccountHash(account), account)
}

func getPlayerKey(p *player.Player) string {
	return fmt.Sprintf("%sPlayer:{%d}:%d",
		config.Instance().App.Prefix, p.Zone, p.Uid)
}

func getShowIdBaseKey(p *player.Player) string {
	return fmt.Sprintf("%sShowIdBase:{%d}",
		config.Instance().App.Prefix, p.Zone)
}

func getShowIdUidMappingKey(p *player.Player) string {
	return fmt.Sprintf("%sShowIdUidMapping:{%d}",
		config.Instance().App.Prefix, p.Zone)
}

func OnLoginReq(a player.Agent, req *message.LoginReq) (
	res message.LoginRes, err error) {
	res = message.LoginRes{}

	conn := dbredis.Conn()
	defer conn.Close()

	// 锁住账号，防止并发登录、创角
	lockKey := getAccountLockKey(req.Account)

	if err = dbredis.Lock(conn, lockKey, 120); err != nil {
		slog.Error("manager: OnLoginReq lock account [%s] error [%s]",
			req.Account, err)
		return
	}
	defer dbredis.Unlock(conn, lockKey)

	accountKey := getAccountKey(req.Account)

	var (
		token     string
		tokenTime int64
		uToken    uuid.UUID
	)

	if req.Token == "" { // sdk登录，可能再传一个SdkToken之类
		if uToken, err = uuid.NewUUID(); err != nil {
			err = fmt.Errorf("manager: OnLoginReq account [%s] new uuid error [%w]",
				req.Account, err)
			return
		}

		token = uToken.String()
	} else { // 重登录
		if token, tokenTime, err = getAccountTokenInfo(conn, req.Account); err != nil {
			err = fmt.Errorf("manager: OnLoginReq getAccountTokenInfo [%s] error [%w]",
				req.Account, err)
			return
		}

		if err = checkTokenInfoValid(req.Token, token, tokenTime); err != nil {
			err = fmt.Errorf("manager: OnLoginReq checkTokenInfoValid [%s] error [%w]",
				req.Account, err)
			return
		}
	}

	tokenTime = time.Now().Unix()

	if _, err = conn.Do("HSET", accountKey, "Token", token, "TokenTime", tokenTime); err != nil {
		return
	}

	a.SetLoginInfo(player.LoginInfo{
		Account: req.Account,
		Token:   token,
	})

	res.Token = token
	return
}

func OnEnterReq(a player.Agent, req *message.EnterReq) (
	p *player.Player, res message.EnterRes, err error) {

	res = message.EnterRes{}

	conn := dbredis.Conn()
	defer conn.Close()

	// TODO 检查zone是否合法（开启状态 etc）
	// TODO 细化zone的检查，先临时处理
	if req.Zone == 0 {
		req.Zone = 1
	}

	// 锁住账号，防止并发登录、创角
	accountLockKey := getAccountLockKey(a.LoginInfo().Account)

	if err = dbredis.Lock(conn, accountLockKey, 120); err != nil {
		slog.Error("OnEnterReq: lock account [%s] error [%s]",
			a.LoginInfo().Account, err)
		return
	}
	defer dbredis.Unlock(conn, accountLockKey)

	// 进游戏时，再次确认自己是最后一次登录
	var (
		token     string
		tokenTime int64
	)

	if token, tokenTime, err = getAccountTokenInfo(conn, a.LoginInfo().Account); err != nil {
		slog.Error("OnEnterReq: get account [%s] token info error [%w]",
			a.LoginInfo().Account, err)
		return
	}

	if err = checkTokenInfoValid(a.LoginInfo().Token, token, tokenTime); err != nil {
		slog.Error("OnEnterReq: check token [%s] error [%w]",
			a.LoginInfo().Token, err)
		return
	}

	var (
		bytes      []byte
		reply      int
		isNew      bool
		account    = &player.ZoneAccount{}
		accountKey = getAccountKey(a.LoginInfo().Account)
	)

	// 获取账号在该区服的信息
	if bytes, err = redis.Bytes(conn.Do("HGET", accountKey, req.Zone)); err != nil {
		if !errors.Is(err, redis.ErrNil) {
			slog.Error("OnEnterReq: get account [%s] zone [%d] info error [%w]",
				a.LoginInfo().Account, req.Zone, err)
			return
		}

		// 创建新号
		account.Account = a.LoginInfo().Account
		account.Uid = uid.Generate()
		account.Zone = req.Zone
		account.CreateTime = time.Now().Unix()

		if bytes, err = account.Marshal(); err != nil {
			return
		}

		if reply, err = redis.Int(conn.Do("HSETNX", accountKey, req.Zone, bytes)); err == nil {
			return
		}

		if reply != 1 {
			slog.Error("OnEnterReq: first set account [%s] zone [%d] is already exist",
				a.LoginInfo().Account, req.Zone)
			return
		} else {
			if err = account.Unmarshal(bytes); err != nil {
				return
			}
		}
	} else {
		if err = account.Unmarshal(bytes); err != nil {
			return
		}
	}

	p = player.New(a, account.Account, account.Uid, account.Zone)
	p.EnterProcessing = true
	defer func() {
		p.EnterProcessing = false
	}()

	if isNew, err = OnPlayerOnline(p); err != nil {
		return
	}

	if err = SavePlayerOnEnter(conn, p, isNew); err != nil {
		return
	}

	if res.PlayerInfo, err = BuildPlayerInfoMsgPack(p); err != nil {
		return
	}

	return
}

func getAccountTokenInfo(conn redis.Conn, account string) (token string, tokenTime int64, err error) {
	var (
		values     []any
		accountKey = getAccountKey(account)
	)

	if values, err = redis.Values(conn.Do("HMGET", accountKey, "Token", "TokenTime")); err != nil {
		return
	}

	token = string(values[0].([]byte))

	if tokenTime, err = strconv.ParseInt(string(values[1].([]byte)), 10, 64); err != nil {
		return
	}

	return
}

func checkTokenInfoValid(reqToken string, token string, tokenTime int64) (err error) {
	if token != reqToken {
		return errTokenMismatch
	}

	if tokenTime+60*60*24 < time.Now().Unix() {
		return errTokenExpired
	}

	return nil
}

func getAccountHash(account string) uint64 {
	return xxhash.Sum64String(account)
}

func getAccountLockKey(account string) string {
	return fmt.Sprintf("%sAccountLock:{%d}:{%s}",
		config.Instance().App.Prefix, getAccountHash(account), account)
}

func OnPlayerOnline(p *player.Player) (isNew bool, err error) {
	if err = p.Load(getPlayerKey(p)); err != nil {
		if !errors.Is(err, redis.ErrNil) {
			return
		}

		p.UpdateUid(p.Uid)
		p.UpdateAccount(p.Account)
		p.UpdateZone(p.Zone)
		p.UpdateCreateTime(time.Now().Unix())

		if err = OnPlayerCreate(p); err != nil {
			return
		}

		isNew = true
	}

	p.UpdateLastLoginTime(time.Now().Unix())
	p.UpdateLastLoginIp(p.RemoteAddr().String())

	return
}

func BuildPlayerInfoMsgPack(p *player.Player) (info message.PlayerInfo, err error) {
	info = message.PlayerInfo{}

	info.Name = p.Name
	info.Uid = p.Uid
	info.Gender = p.Gender

	if info.Name == "" {
		info.Name = fmt.Sprintf("%d", info.Uid)
	}

	return
}
