package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
	_ "github.com/lib/pq"
)

type PostgreStorage struct {
	db *sql.DB
}

func getDBConnURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRE_HOST"), os.Getenv("POSTGRE_PORT"), os.Getenv("POSTGRE_USER"), os.Getenv("POSTGRE_PASS"), os.Getenv("POSTGRE_DB"))
}

func New() *PostgreStorage {
	postgreURL := getDBConnURL()

	conn, err := sql.Open("postgres", postgreURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB connected")

	return &PostgreStorage{db: conn}
}

func (postgre *PostgreStorage) AddNewItem(newItem itemsModel.Item) error {
	_, err := postgre.db.Exec(
		`INSERT INTO "items"(name, price, number, description) VALUES ($1, $2, $3, $4) RETURNING id`,
		newItem.Name,
		newItem.Price,
		newItem.ItemsNumber,
		newItem.Description,
	)

	if err != nil {
		return err
	}

	return nil
}

func (postgre *PostgreStorage) GetAllItems() ([]byte, error) {
	var itemsSlice []itemsModel.Item

	rows, err := postgre.db.Query(`SELECT name, price, number, description FROM "items"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item itemsModel.Item
		if err := rows.Scan(&item.Name, &item.Price, &item.ItemsNumber, &item.Description); err != nil {
			return nil, err
		}
		itemsSlice = append(itemsSlice, item)
	}

	res, err := json.Marshal(itemsSlice)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (postgre *PostgreStorage) GetItem(name string) ([]byte, error) {
	var item itemsModel.Item
	err := postgre.db.QueryRow(`SELECT name, price, number, description FROM "items" WHERE name=$1`, name).Scan(&item.Name, &item.Price, &item.ItemsNumber, &item.Description)
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (postgre *PostgreStorage) DeleteItem(itemName string) error {
	_, err := postgre.db.Exec(
		`DELETE FROM "items" WHERE name=&1`,
		itemName,
	)

	if err != nil {
		return err
	}

	return nil
}

func (postgre *PostgreStorage) UpdateItem(itemName string) ([]byte, error) {
	var item itemsModel.Item

	err := postgre.db.QueryRow(
		`UPDATE "user" SET number=number - 1 WHERE name=$1 RETURNING name, price, number, description`,
		itemName,
	).Scan(&item.Name, &item.Price, &item.ItemsNumber, &item.Description)
	if err != nil {
		return nil, err
	}

	// return result
	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}
