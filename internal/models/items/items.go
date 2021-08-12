package items

type Item struct {
	ID          int    `json:"id,omitempty" csv:"name"`
	Name        string `json:"name,omitempty" csv:"name"`
	Price       int32  `json:"price,omitempty" csv:"price"`
	ItemsNumber int32  `json:"itemsNumber,omitempty" csv:"itemsNumber"`
	Description string `json:"desc,omitempty" csv:"desc"`
}

// all main services for Item to work with DB
type ItemStorageServices interface {
	AddNewItem(Item) (int, error)
	GetAllItems() ([]byte, error)
	GetItem(int) ([]byte, error)
	DeleteItem(int) (bool, error)
	UpdateItem(int) ([]byte, error)
}
