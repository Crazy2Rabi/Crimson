package genMessage

type Packets struct {
	Packets []Packet `xml:"packet"`
	Refs    []Ref    `xml:"ref"`
}

// 消息
type Packet struct {
	Name   string  `xml:"name,attr"`
	Id     string  `xml:"id,attr"`
	Desc   string  `xml:"des,attr"`
	Fields []Field `xml:"field"`
}

// 结构体
type Ref struct {
	Name   string  `xml:"name,attr"`
	Desc   string  `xml:"des,attr"`
	Fields []Field `xml:"field"`
}

type Field struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	RefType string `xml:"refType,attr"`
	Desc    string `xml:"des,attr"`
}
