package redis

import (
	"context"
	"encoding/json"
	"strconv"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
	redis "github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func New() *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisStorage{client: rdb}
}

// search not exist id to add to db
// recommend add len of db items as first argument
func (db *RedisStorage) getNewID(id int) (int, error) {
	_, err := db.client.Get(context.Background(), strconv.Itoa(id)).Result()
	switch {
	case err == redis.Nil:
		return id, nil
	case err != nil:
		return 0, err
	}

	return db.getNewID(id + 1)
}

func (db *RedisStorage) AddNewItem(newItem itemsModel.Item) (int, error) {
	lastID, err := db.client.DBSize(context.Background()).Result()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}

		lastID = 0 // start init if not exist
	}

	id := int(lastID) + 1 // use

	newID, err := db.getNewID(id)
	if err != nil {
		return 0, err
	}
	_, err = db.client.Set(context.Background(), strconv.Itoa(newID), &newItem, 0).Result()
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (db *RedisStorage) GetAllItems() ([]byte, error) {
	var itemsSlice []itemsModel.Item

	iter := db.client.Scan(context.Background(), 0, "*", 0).Iterator()
	for iter.Next(context.Background()) {
		res, err := db.client.Get(context.Background(), iter.Val()).Bytes()
		if err != nil {
			return nil, err
		}

		var item itemsModel.Item
		err = json.Unmarshal(res, &item)
		if err != nil {
			return nil, err
		}

		// add id to result
		id, err := strconv.Atoi(iter.Val())
		if err != nil {
			return nil, err
		}
		item.ID = id

		itemsSlice = append(itemsSlice, item)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	res, err := json.Marshal(itemsSlice)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (db *RedisStorage) GetItem(id int) ([]byte, error) {
	res, err := db.client.Get(context.Background(), strconv.Itoa(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

func (db *RedisStorage) DeleteItem(id int) (bool, error) {
	rows, err := db.client.Del(context.Background(), strconv.Itoa(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	if rows <= 0 {
		return false, nil
	}

	return true, nil
}

func (db *RedisStorage) UpdateItem(id int) ([]byte, error) {
	foundItem, err := db.client.Get(context.Background(), strconv.Itoa(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var item itemsModel.Item
	if err = json.Unmarshal(foundItem, &item); err != nil {
		return nil, err
	}

	// update value count
	item.ItemsNumber--
	if item.ItemsNumber < 0 {
		_, err := db.client.Del(context.Background(), strconv.Itoa(id)).Result()
		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		return nil, nil

	}

	_, err = db.client.Set(context.Background(), strconv.Itoa(id), &item, 0).Result()
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}
