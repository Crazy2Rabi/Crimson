package main

import (
	"Common/GenModule/genOthers/genJSON"
)

func main() {
	// 生成JSON
	if err := genJSON.GenJSON(); err != nil {
		panic(err)
	}

}
