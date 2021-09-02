package redis

import (
	"encoding/json"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
)

type itemRedis itemsModel.Item

func (i itemRedis) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
