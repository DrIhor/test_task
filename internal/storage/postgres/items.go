package postgres

import (
	"context"
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

func (postgre *PostgreStorage) AddNewItem(ctx context.Context, newItem itemsModel.Item) (int, error) {
	var newItemID int

	// DB request
	err := postgre.db.QueryRowContext(
		ctx,
		`INSERT INTO "items"(name, price, number, description) VALUES ($1, $2, $3, $4) RETURNING id`,
		newItem.Name,
		newItem.Price,
		newItem.ItemsNumber,
		newItem.Description,
	).Scan(&newItemID)

	if err != nil {
		return 0, err
	}

	return newItemID, nil
}

func (postgre *PostgreStorage) GetAllItems(ctx context.Context) ([]byte, error) {
	var itemsSlice []itemsModel.Item

	// DB request
	rows, err := postgre.db.QueryContext(ctx, `SELECT id, name, price, number, description FROM "items"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item itemsModel.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.ItemsNumber, &item.Description); err != nil {
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

func (postgre *PostgreStorage) GetItem(ctx context.Context, id int) ([]byte, error) {
	var item itemsModel.Item

	err := postgre.db.QueryRowContext(
		ctx,
		`SELECT name, price, number, description FROM "items" WHERE id=$1`, id).Scan(&item.Name, &item.Price, &item.ItemsNumber, &item.Description)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (postgre *PostgreStorage) DeleteItem(ctx context.Context, id int) (bool, error) {

	res, err := postgre.db.ExecContext(
		ctx,
		`DELETE FROM "items" WHERE id=$1`,
		id,
	)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected <= 0 {
		return false, nil
	}

	return true, nil
}

func (postgre *PostgreStorage) UpdateItem(ctx context.Context, id int) ([]byte, error) {
	var item itemsModel.Item

	err := postgre.db.QueryRowContext(
		ctx,
		`UPDATE "items" SET number=number - 1 WHERE id=$1 RETURNING name, price, number, description`,
		id,
	).Scan(&item.Name, &item.Price, &item.ItemsNumber, &item.Description)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if item.ItemsNumber < 0 {
		res, err := postgre.db.ExecContext(
			ctx,
			`DELETE FROM "items" WHERE id=$1`,
			id,
		)
		if err != nil {
			return nil, err
		}

		_, err = res.RowsAffected()
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// return result
	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}
