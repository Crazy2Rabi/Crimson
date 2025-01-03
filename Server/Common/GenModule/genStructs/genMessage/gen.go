package genMessage

import (
	"Common/Framework/config"
	"Common/GenModule/GenFile"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

func GenMessage() error {
	startTime := time.Now()

	// 相对于build.bat中的路径
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// fmt.Println("当前目录：", dir)
	messagePath := filepath.Join(dir, config.Instance().GenModuleConfig.MessagePath)
	genMessagePath := filepath.Join(dir, config.Instance().GenModuleConfig.MessageGenPath)
	genStructPath := filepath.Join(dir, config.Instance().GenModuleConfig.MessageDefPath)

	data, err := os.ReadFile(messagePath)
	if err != nil {
		return err
	}

	var packets Packets

	if err = xml.Unmarshal(data, &packets); err != nil {
		return err
	}

	// 生成通用结构体
	path := filepath.Join(genStructPath, "ref.go")
	GenFile.GenFile(path, RefText, packets, template.FuncMap{})

	// 生成消息
	// 首字母大写，否则修饰符是private
	/*for _, p := range packets.Packets {
		p.Name = strings.Title(p.Name)
		for _, f := range p.Fields {
			f.Name = strings.Title(f.Name)
		}
	}
	*/

	// 转换类型
	for j, p := range packets.Packets {
		for k, f := range p.Fields {
			if t, ok := GenFile.TypeConverter[f.Type]; !ok {
				packets.Packets[j].Fields[k].Type = f.Type
			} else {
				packets.Packets[j].Fields[k].Type = t
			}
		}
	}

	path = filepath.Join(genMessagePath, "message.go")
	GenFile.GenFile(path, ModelText, packets, template.FuncMap{})

	fmt.Printf("gen message ok! Total time cost %v\n\n", time.Since(startTime))
	return nil
}
