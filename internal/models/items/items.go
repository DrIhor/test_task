package items

type Item struct {
	Name        string `json:"name"`
	Price       int32  `json:"price"`
	ItemsNumber int32  `json:"itemsNumber"`
	Description string `json:"desc"`
}

// all main services for Item to work with DB
type ItemStorageServices interface {
	AddNewItem(Item) error
	GetAllItems() ([]byte, error)
	GetItem(string) ([]byte, error)
	DeleteItem(string) error
	UpdateItem(string) ([]byte, error)
}
