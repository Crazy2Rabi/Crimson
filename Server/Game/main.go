package main

import (
	"Common/Framework"
	"Common/Framework/codec"
	"Common/Framework/dbredis"
	"Common/Framework/tableConfig"
	"Common/Framework/uid"
	"Game/services/agent"
	"Game/services/server"
)

func main() {
	Framework.Run(Framework.WithInit(uid.Init, codec.Init, dbredis.Init, tableConfig.Init),
		Framework.WithServices(agent.New, server.New),
	)
}
