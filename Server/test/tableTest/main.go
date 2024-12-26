package main

import (
	"Common/Framework/tableConfig"
	"bufio"
	"fmt"
	"os"
)

func main() {
	var err error
	err = tableConfig.Init()
	if err != nil {
		fmt.Println(err)
	}

	Test1()
}

// 重载测试
func Test1() {
	t := tableConfig.Instance().TbConfigItem.FindById(1002)
	if t != nil {
		t.Prize.PrintInfo()
	}

	// 重载
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请修改文件，按任意键继续")
	_, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tableConfig.Reload()
	if err != nil {
		fmt.Println(err)
	}
	t = tableConfig.Instance().TbConfigItem.FindById(1002)
	if t != nil {
		t.Prize.PrintInfo()
	}
}
