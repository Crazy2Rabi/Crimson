package uid

import (
	"Common/Framework/config"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"log/slog"
	"math"
)

var node *snowflake.Node

func Init() (err error) {
	if config.Instance().App.Id > math.MaxInt64 {
		err = fmt.Errorf("uid overflow")
		return
	}

	node, err = snowflake.NewNode(int64(config.Instance().App.Id))
	slog.Info("uid init ok")
	return
}

func Generate() uint64 {
	return uint64(node.Generate())
}
