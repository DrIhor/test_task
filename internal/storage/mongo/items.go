package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"

	"github.com/google/uuid"
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

func New() (*MongoStorage, error) {
	mongoURL := getDBConnURL()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &MongoStorage{db: client.Database("shop")}, err
}

func (mg *MongoStorage) AddNewItem(ctx context.Context, newItem itemsModel.Item) (string, error) {
	collection := mg.db.Collection("items")

	// set new id
	newItem.ID = uuid.New().String()

	res, err := collection.InsertOne(ctx, newItem)
	if err != nil {
		return "", err
	}

	strID := fmt.Sprintf("%v", res.InsertedID)

	return strID, nil
}

func (mg *MongoStorage) GetAllItems(ctx context.Context) ([]byte, error) {
	collection := mg.db.Collection("items")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var items []itemsModel.Item
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}

	res, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (mg *MongoStorage) GetItem(ctx context.Context, id string) ([]byte, error) {
	collection := mg.db.Collection("items")

	var item itemsModel.Item
	err := collection.FindOne(ctx, bson.M{
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

func (mg *MongoStorage) DeleteItem(ctx context.Context, id string) (bool, error) {
	collection := mg.db.Collection("items")

	res, err := collection.DeleteOne(ctx, bson.M{
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

func (mg *MongoStorage) UpdateItem(ctx context.Context, id string) ([]byte, error) {
	collection := mg.db.Collection("items")

	var item itemsModel.Item

	after := options.After
	err := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{
				"itemsNumber": -1,
			},
		},
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&item)
	if err != nil {
		return nil, err
	}

	if item.ItemsNumber < 0 {
		_, err := collection.DeleteOne(
			ctx,
			bson.M{
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
