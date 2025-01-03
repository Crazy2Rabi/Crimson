package config

import (
	"github.com/jinzhu/configor"
	"strconv"
	"sync"
)

var config Configure
var once sync.Once

type Configure struct {
	App appConfig

	Redis           redisConfig
	GenModuleConfig genModuleConfig
	TableConfig     tableConfig
}

type appConfig struct {
	Id     uint64 `default:"1"` // app id
	Zone   uint64 `default:"1"` // app zone id
	Name   string `default:"Crimson-Server"`
	Addr   string `default:":9000"`      // websocket服务器监听端口
	Prefix string `default:"FengShuang"` // 账号前缀
}

type genModuleConfig struct {
	MessagePath    string `default:"../../Message/packet.xml"` // xml消息目录
	MessageGenPath string `default:"../Common/message"`        // go消息生成目录
	MessageDefPath string `default:"../Common/message"`        // 消息中的结构体生成目录

	TablesPath       string `default:"../../Table"`          // 表格目录
	TablesStructPath string `default:"../Common/Table"`      // 表结构生成目录
	TablesJSONPath   string `default:"../Common/Table/json"` // 表生成JSON目录
}

type tableConfig struct {
	Path string `default:"../Common/Table/json/"` // 表json目录
}

type redisConfig struct {
	Addr      string `default:"localhost:6379"`
	DB        int    `default:"0"`
	MaxIdle   int    `default:"10"`
	MaxActive int    `default:"20"`
	Username  string `default:"root"`
	Password  string `default:""`
}

func (c *Configure) init() (err error) {
	if err = configor.Load(c, "../config.json"); err != nil {
		return err
	}

	c.App.Prefix += strconv.FormatUint(c.App.Zone, 10)

	return
}

func Instance() *Configure {
	once.Do(func() {
		config = Configure{}
		if err := config.init(); err != nil {
			panic(err)
		}
	})
	return &config
}
