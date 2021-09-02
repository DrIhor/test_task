package messages

type responce struct {
	Msg string `json:"msg"`
}

type Message interface {
	CreateMsgResp(string) ([]byte, error)
}
