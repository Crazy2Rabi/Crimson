package codec

import (
	"Common/message"
	"encoding/binary"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"log/slog"
	"reflect"
	"strconv"
)

var codec *Codec

const cmdFieldSize = 2 // 协议号占用空间

type Codec struct {
	cmdTypes map[uint16]reflect.Type
	typeCmds map[reflect.Type]uint16
}

func NewCodec() *Codec {
	return &Codec{
		cmdTypes: make(map[uint16]reflect.Type),
		typeCmds: make(map[reflect.Type]uint16),
	}
}

func (c *Codec) registerAllMessage(m interface{}) error {
	t := reflect.TypeOf(m)

	for i := 0; i < t.NumField(); i++ {
		cmdId, err := strconv.ParseUint(t.Field(i).Tag.Get("Id"), 10, 32)
		if err != nil {
			return err
		}

		cmd := uint16(cmdId)
		if _, ok := c.cmdTypes[cmd]; ok {
			return fmt.Errorf("codec: register message [%v] [%v] is duplicated", cmd, t.Field(i).Type)
		}

		c.cmdTypes[cmd] = t.Field(i).Type
		c.typeCmds[t.Field(i).Type] = cmd
	}

	return nil
}

func (c *Codec) messageType(cmd uint16) (t reflect.Type, err error) {
	if t, ok := c.cmdTypes[cmd]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("message cmd [%d] is not register", cmd)
}

func (c *Codec) messageCmd(t reflect.Type) (cmd uint16, err error) {
	if cmd, ok := c.typeCmds[t]; ok {
		return cmd, nil
	}
	return 0, fmt.Errorf("message type [%v] is not register", t)
}

func (c *Codec) encode(m interface{}) (data []byte, err error) {
	var cmd uint16

	if cmd, err = c.messageCmd(reflect.TypeOf(m)); err != nil {
		err = fmt.Errorf("encode message [%v] message cmd error [%w]", m, err)
		return
	}

	var buf []byte
	if buf, err = msgpack.Marshal(m); err != nil {
		err = fmt.Errorf("encode cmd [%#04x] Marshal error [%w]", cmd, err)
		return
	}

	dataSize := cmdFieldSize + len(buf)
	data = make([]byte, dataSize)
	binary.BigEndian.PutUint16(data, cmd)

	copy(data[cmdFieldSize:], buf)
	return
}

func (c *Codec) decode(data []byte) (m interface{}, err error) {
	headerSize := cmdFieldSize
	if len(data) < headerSize {
		err = fmt.Errorf("decode data size [%d] < header size [%d]", len(data), headerSize)
		return
	}

	cmd := binary.BigEndian.Uint16(data)

	var t reflect.Type

	if t, err = c.messageType(cmd); err != nil {
		err = fmt.Errorf("decode cmd [%d] message type error [%w]", cmd, err)
		return
	}

	p := reflect.New(t)

	if err = msgpack.Unmarshal(data[headerSize:], p.Interface()); err != nil {
		err = fmt.Errorf("decode message cmd [%d] Unmarshal error [%w]", cmd, err)
		return
	}

	m = p.Elem().Interface()
	return
}

func Init() error {
	codec = NewCodec()
	if err := codec.registerAllMessage(message.Message{}); err != nil {
		return err
	}
	slog.Info("codec init ok")
	return nil
}

func MessageType(cmd uint16) (t reflect.Type, err error) {
	return codec.messageType(cmd)
}

func MessageCmd(t reflect.Type) (cmd uint16, err error) {
	return codec.messageCmd(t)
}

func Encode(m interface{}) (data []byte, err error) {
	return codec.encode(m)
}

func Decode(data []byte) (m interface{}, err error) {
	return codec.decode(data)
}
