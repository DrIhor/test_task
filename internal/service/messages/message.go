package messages

import (
	"encoding/json"
)

func NewMessage() *responce {
	return &responce{}
}

func (ms *responce) CreateMsgResp(message string) ([]byte, error) {

	ms.Msg = message

	return json.Marshal(ms)
}
