// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Item struct {
	ID          *int    `json:"id"`
	Name        string  `json:"name"`
	Price       int     `json:"price"`
	ItemsNumber int     `json:"itemsNumber"`
	Desc        *string `json:"desc"`
}

type Iteminput struct {
	Name        string  `json:"name"`
	Price       int     `json:"price"`
	ItemsNumber int     `json:"itemsNumber"`
	Desc        *string `json:"desc"`
}