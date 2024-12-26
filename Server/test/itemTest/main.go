package main

import (
	"Common/Table"
	"Common/Utils/def"
	"Game/manager"
	"Game/player"
	"fmt"
)

func main() {
	/*Framework.Run(Framework.WithInit(uid.Init, codec.Init, dbredis.Init),
		Framework.WithServices(agent.New, server.New),
	)*/

	var err error
	var p player.Player = player.Player{
		Player: &def.Player{},
	}

	err = Table.Table.Init()
	if err != nil {
		fmt.Println(err)
	}

	items := def.Items{
		def.NewItem(1003, 11),
		def.NewItem(1001, 10000),
		def.NewItem(1003, 2),
		def.NewItem(1003, 9),
		def.NewItem(1001, 10000),
		def.NewItem(1002, 1000),
	}

	fmt.Println("\n增加物品:")
	err = manager.AddItems(&p, &items)
	if err != nil {
		fmt.Println(err)
	}
	manager.PrintItemsInfo(p)

	fmt.Println("\n消耗物品:")
	items = def.Items{
		def.NewItem(1001, 3000),
		def.NewItem(1002, 200),
		def.NewItem(1003, 21),
	}
	err = manager.SubItems(&p, &items)
	if err != nil {
		fmt.Println(err)
	}
	manager.PrintItemsInfo(p)
}
