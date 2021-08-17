package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoStorage struct {
	db *mongo.Database
}

func getDBConnURL() string {
	return fmt.Sprintf("mongodb://%s", os.Getenv("MONGO_ADDR"))
}

func New() *MongoStorage {
	mongoURL := getDBConnURL()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connected")

	return &MongoStorage{db: client.Database("shop")}
}

func (mg *MongoStorage) AddNewItem(newItem itemsModel.Item) (int, error) {
	collection := mg.db.Collection("items")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// check id of last record
	var lastrecord itemsModel.Item
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	if err := collection.FindOne(ctx, bson.M{}, opts).Decode(&lastrecord); err != nil {
		if err != mongo.ErrNoDocuments {
			return 0, err
		}

		lastrecord.ID = 0
	}

	// set new id
	newItem.ID = lastrecord.ID + 1

	res, err := collection.InsertOne(ctx, newItem)
	if err != nil {
		return 0, err
	}

	strID := fmt.Sprintf("%v", res.InsertedID)
	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (mg *MongoStorage) GetAllItems() ([]byte, error) {
	collection := mg.db.Collection("items")
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var items []itemsModel.Item
	if err := cur.All(context.Background(), &items); err != nil {
		return nil, err
	}

	res, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mg *MongoStorage) GetItem(id int) ([]byte, error) {
	collection := mg.db.Collection("items")

	var item itemsModel.Item
	err := collection.FindOne(context.Background(), bson.M{
		"_id": id,
	}).Decode(&item)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mg *MongoStorage) DeleteItem(id int) (bool, error) {
	collection := mg.db.Collection("items")

	res, err := collection.DeleteOne(context.Background(), bson.M{
		"_id": id,
	})
	if err != nil {
		return false, err
	}

	if res.DeletedCount <= 0 {
		return false, nil
	}

	return true, nil

}

func (mg *MongoStorage) UpdateItem(id int) ([]byte, error) {
	collection := mg.db.Collection("items")

	var item itemsModel.Item
	err := collection.FindOneAndUpdate(context.Background(),
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{
				"itemsNumber": -1,
			},
		}).Decode(&item)
	if err != nil {
		return nil, err
	}

	if item.ItemsNumber <= 0 {
		_, err := collection.DeleteOne(context.Background(), bson.M{
			"_id": id,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}
