package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Redis is connected")

	return &RedisStorage{client: rdb}
}

// search not exist id to add to db
// recommend add len of db items as first argument
func (db *RedisStorage) getNewID(ctx context.Context, id int) (int, error) {
	_, err := db.client.Get(ctx, strconv.Itoa(id)).Result()
	switch {
	case err == redis.Nil:
		return id, nil
	case err != nil:
		return 0, err
	}

	return db.getNewID(ctx, id+1)
}

func (db *RedisStorage) AddNewItem(ctx context.Context, newItem itemsModel.Item) (int, error) {
	lastID, err := db.client.DBSize(ctx).Result()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}

		lastID = 0 // start init if not exist
	}

	id := int(lastID) + 1 // use

	newID, err := db.getNewID(ctx, id)
	if err != nil {
		return 0, err
	}
	_, err = db.client.Set(ctx, strconv.Itoa(newID), &newItem, 0).Result()
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (db *RedisStorage) GetAllItems(ctx context.Context) ([]byte, error) {
	var itemsSlice []itemsModel.Item

	iter := db.client.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		res, err := db.client.Get(ctx, iter.Val()).Bytes()
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

func (db *RedisStorage) GetItem(ctx context.Context, id int) ([]byte, error) {
	res, err := db.client.Get(ctx, strconv.Itoa(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

func (db *RedisStorage) DeleteItem(ctx context.Context, id int) (bool, error) {
	rows, err := db.client.Del(ctx, strconv.Itoa(id)).Result()
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

func (db *RedisStorage) UpdateItem(ctx context.Context, id int) ([]byte, error) {
	foundItem, err := db.client.Get(ctx, strconv.Itoa(id)).Bytes()
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
		_, err := db.client.Del(ctx, strconv.Itoa(id)).Result()
		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		return nil, nil

	}

	_, err = db.client.Set(ctx, strconv.Itoa(id), &item, 0).Result()
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}
