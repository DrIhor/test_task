package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisStorage struct {
	client *redis.Client
}

func New() (*RedisStorage, error) {

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       db,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStorage{client: rdb}, nil
}

func (db *RedisStorage) AddNewItem(ctx context.Context, newItem itemsModel.Item) (string, error) {
	newID := uuid.New()

	newItemRd := itemRedis(newItem)

	_, err := db.client.Set(ctx, newID.String(), &newItemRd, 0).Result()
	if err != nil {
		return "", err
	}

	return newID.String(), nil
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
		item.ID = iter.Val()

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

func (db *RedisStorage) GetItem(ctx context.Context, id string) ([]byte, error) {
	res, err := db.client.Get(ctx, id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

func (db *RedisStorage) DeleteItem(ctx context.Context, id string) (bool, error) {
	rows, err := db.client.Del(ctx, id).Result()
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

func (db *RedisStorage) UpdateItem(ctx context.Context, id string) ([]byte, error) {
	foundItem, err := db.client.Get(ctx, id).Bytes()
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
		_, err := db.client.Del(ctx, id).Result()
		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		return nil, nil

	}

	itemRd := itemRedis(item)
	_, err = db.client.Set(ctx, id, &itemRd, 0).Result()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}
