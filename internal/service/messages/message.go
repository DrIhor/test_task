package messages

import (
	"encoding/json"

	messageModel "github.com/DrIhor/test_task/internal/models/messages"
)

func CreateMsgResp(message string) []byte {

	dataResp := messageModel.Responce{
		Msg: message,
	}

	res, _ := json.Marshal(dataResp)
	return res

}
