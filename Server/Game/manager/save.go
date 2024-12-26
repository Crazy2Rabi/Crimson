package manager

import (
	"Common/Utils/def"
	"Game/player"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func checkUpdatedFV(fvs []def.FV, isNew bool) (err error) {
	newPlayerKeys := map[string]int{
		"Uid":        -1,
		"Account":    -1,
		"Zone":       -1,
		"CreateTime": -1,
	}

	for _, v := range fvs {
		if _, ok := newPlayerKeys[v.Field]; ok {
			newPlayerKeys[v.Field]++
		}
	}

	for _, v := range newPlayerKeys {
		if isNew {
			if v == -1 {
				err = errors.New("manager: create player lost data")
				return
			}
		} else {
			if v == 0 {
				err = errors.New("manager: save player illegal data")
				return
			}
		}
	}

	return
}

func SavePlayerOnEnter(conn redis.Conn, p *player.Player, isNew bool) (err error) {
	return p.Process(func(fvs []def.FV) (err error) {
		if err = checkUpdatedFV(fvs, isNew); err != nil {
			return
		}

		// 放在后面判断，创角时先进行必要字段检查
		if len(fvs) == 0 {
			return
		}

		p.Temp = p.Temp[:0]

		playerKey := getPlayerKey(p)
		showIdBaseKey := getShowIdBaseKey(p)
		showIdUidMappingKey := getShowIdUidMappingKey(p)

		// KEYS
		p.Temp = append(p.Temp, playerKey, showIdBaseKey, showIdUidMappingKey)

		// ARGV
		var (
			isNewValue       = 0
			minShowId  int32 = 10000000
			maxShowId  int32 = 100000000
		)

		if isNew {
			isNewValue = 1
		}

		p.Temp = append(p.Temp,
			isNewValue,
			p.Uid,
			p.LoginInfo().Token,
			time.Now().Unix(),
			minShowId,
			maxShowId,
		)

		for _, v := range fvs {
			p.Temp = append(p.Temp, v.Field)
			p.Temp = append(p.Temp, v.Value)
		}

		script := redis.NewScript(4, `
			-- 参数整理
			local player_key = KEYS[1]
			local show_id_base_key = KEYS[2]
			local show_id_uid_mapping_key = KEYS[3]
			local show_uid_id_mapping_key = KEYS[4]

			local is_new_value = tonumber(ARGV[1])
			local uid = ARGV[2]
			local token = ARGV[3]
			local token_time = ARGV[4]
			local min_show_id = tonumber(ARGV[5])
			local max_show_id = tonumber(ARGV[6])

			local player_fv = {}
			for i = 7, #ARGV, 2 do
				table.insert(player_fv, ARGV[i])
				table.insert(player_fv, ARGV[i+1])
			end

			-- 检查player key的是否存在
			local reply = redis.call('EXISTS', player_key)
			if is_new_value == 1 then
				-- 新号必须没有player
				if reply ~= 0 then
					return -1
				end
			else
				-- 老号必须有player
				if reply ~= 1 then
					return -2
				end
			end

			-- 新号自增生成本服8位show id
			local show_id = 0

			if is_new_value == 1 then
				show_id = redis.call('INCR', show_id_base_key)
				show_id = show_id + min_show_id
				
				if show_id >= max_show_id then
					return -3
				end

				-- 保存本服show id和uid的映射，方便查找
				local reply2 = redis.call('HSETNX', show_id_uid_mapping_key, show_id, uid)
				if reply2 ~= 1 then
					return -4
				end
			end

			redis.call('HSET', player_key, unpack(player_fv))

			-- 新号写玩家的show id
			if is_new_value == 1 and show_id > 0 then
				redis.call('HSET', player_key, 'ShowId', show_id)
			end

			-- 登录强制写Token和TokenTime
			redis.call('HSET', player_key, 'Token', token, 'TokenTime', token_time)

			return 0
		`)

		var reply int

		if reply, err = redis.Int(script.Do(conn, p.Temp...)); err != nil {
			return
		}

		if reply != 0 {
			err = fmt.Errorf("player: uid [%d] enter save reply [%d]",
				p.Uid, reply)
			return
		}

		if isNew {
			var showId int

			if showId, err = redis.Int(conn.Do("HGET", playerKey, "showId")); err != nil {
				return
			}

			p.ShowId = int32(showId)

			if p.ShowId <= minShowId || p.ShowId >= maxShowId {
				err = errors.New("manager: showId overflow")
				return
			}
		}

		return
	})
}
