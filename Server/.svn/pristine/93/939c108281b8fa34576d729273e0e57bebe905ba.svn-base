package main

import (
	"Common/GenModule/genStructs/genMessage"
	"Common/GenModule/genStructs/genTable"
)

func main() {
	// 生成消息
	err := genMessage.GenMessage()
	if err != nil {
		panic(err)
	}

	// 生成表结构
	err = genTable.LoadTableConfigs()
	if err != nil {
		panic(err)
	}
}
