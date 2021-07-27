package iteams

type Iteam struct {
	Name         string `json:"name"`
	Price        int    `json:"price"`
	IteamsNumber int    `json:"iteamsNumber"`
	Description  string `json:"desc"`
}

type IteamServices interface {
	AddNewIteam(Iteam) error
	GetAllIteams() ([]byte, error)
	GetIteam(string) ([]byte, error)
	DeleteIteam(string) error
	UpdateIteam(string) ([]byte, error)
}
