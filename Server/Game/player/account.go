package player

import "encoding/json"

type ZoneAccount struct {
	Account    string
	Uid        uint64
	Zone       uint64
	Channel    uint64
	CreateTime int64
}

func (za ZoneAccount) Marshal() (data []byte, err error) {
	data, err = json.Marshal(za)
	return
}

func (za *ZoneAccount) Unmarshal(data []byte) (err error) {
	err = json.Unmarshal(data, za)
	return
}
