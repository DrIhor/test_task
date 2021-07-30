package iteams

type Iteam struct {
	Name         string `json:"name"`
	Price        int    `json:"price"`
	IteamsNumber int    `json:"iteamsNumber"`
	Description  string `json:"desc"`
}

// all main services for Iteam to work with DB
type IteamStorageServices interface {
	AddNewIteam(Iteam) error
	GetAllIteams() ([]byte, error)
	GetIteam(string) ([]byte, error)
	DeleteIteam(string) error
	UpdateIteam(string) ([]byte, error)
}
