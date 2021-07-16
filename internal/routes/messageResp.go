package routes

import "encoding/json"

type responce struct {
	Msg string `json:"msg"`
}

func createMsgResp(message string) []byte {
	dataResp := responce{
		Msg: message,
	}

	res, _ := json.Marshal(dataResp)
	return res

}
