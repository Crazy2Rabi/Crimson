package GenFile

import (
	"bytes"
	"go/format"
	"html/template"
	"os"
	"path/filepath"
	"runtime/debug"
)

func GenFile(path string, tmpl string, data any, fm template.FuncMap) {
	var (
		t    *template.Template
		err  error
		buf  bytes.Buffer
		code []byte
	)

	if t, err = template.New("gen_file").Funcs(fm).Parse(tmpl); err != nil {
		debug.PrintStack()
		panic(err)
	}

	if err = t.Execute(&buf, data); err != nil {
		debug.PrintStack()
		panic(err)
	}

	if code, err = format.Source(buf.Bytes()); err != nil {
		debug.PrintStack()
		panic(err)
	}

	// 生成中间的文件夹
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0755); err != nil {
		debug.PrintStack()
		panic(err)
	}

	if err = os.WriteFile(path, code, 0644); err != nil {
		debug.PrintStack()
		panic(err)
	}
}

// 转换消息或表TYPE行的一些类型
// 如果表中类型为自定义类型，定义在def中，再让策划填上对应的类型
// 不在该容器中的，统一当成自定义类型，生成文件时会修改成 def.类型
var TypeConverter map[string]string = map[string]string{
	"bool":   "bool",
	"int":    "int32",
	"uint32": "uint32",
	"int32":  "int32",
	"uint64": "uint64",
	"int64":  "int64",
	"string": "string",
}
