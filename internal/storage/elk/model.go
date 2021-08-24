package elk

import (
	itemsModel "github.com/DrIhor/test_task/internal/models/items"
)

type SearchELKItems struct {
	Hits struct {
		Hits []struct {
			ID     string          `json:"_id"`
			Index  string          `json:"_index"`
			Score  float64         `json:"_score"`
			Source itemsModel.Item `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type GetELKItems struct {
	Source itemsModel.Item `json:"_source"`
	Found  bool            `json:"found"`
}

type UpdateItemsQuery struct {
	Script struct {
		Source string `json:"source"`
		Lang   string `json:"lang"`
	} `json:"script"`
}
