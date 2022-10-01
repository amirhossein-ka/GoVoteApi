package mongo

import (
	"GoVoteApi/models"
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection = "users"

var (
	ErrNoDocuments error = errors.New("no document found with given id")
	ErrInvalidID   error = errors.New("id is invalid")
)

func getID(result interface{}) string {
	var res string
	if x, ok := result.(primitive.ObjectID); ok {
		res = x.Hex()
	} else {
		res = result.(string)
	}

	return res
}

// CreateUser creates a user document from u and return it's id as string
func (m *mongoImpl) CreateUser(ctx context.Context, u *models.User) (string, error) {
	ctx, cancel := newMongoContext(ctx, m.cfg.Timeout)
	defer cancel()

	res, err := m.userCol.InsertOne(ctx, u)
	if err != nil {
		return "", err
	}
	m.userCol.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: u.UserName}},
			Options: options.Index().SetUnique(true),
		},
	)

	id := getID(res.InsertedID)

	return id, nil
}

func (m *mongoImpl) GetUser(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := newMongoContext(ctx, m.cfg.Timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	var result models.User
	if err := m.userCol.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &result, nil
}

func (m *mongoImpl) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, cancel := newMongoContext(ctx, m.cfg.Timeout)
	defer cancel()
	filter := bson.D{{Key: "username", Value: username}}
	var result models.User
	if err := m.db.Collection(userCollection).FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoDocuments
		}
		return nil, err
	}
	return &result, nil
}

func (m *mongoImpl) UpdateUser(ctx context.Context, u *models.User) error {
	ctx, cancel := newMongoContext(ctx, m.cfg.Timeout)
	defer cancel()

	update := bson.D{{
		Key:   "$set",
		Value: getUpdateArray(u),
	}}
	_, err := m.userCol.UpdateByID(ctx, u.ID, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoDocuments
		}
		return err
	}

	return nil
}

func (m *mongoImpl) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := newMongoContext(ctx, m.cfg.Timeout)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	res, err := m.userCol.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("no document deleted")
	}

	return nil
}

func getUpdateArray(t any) []bson.E {
	var updates []bson.E
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)
		key, ok := v.Type().Field(i).Tag.Lookup("bson")
		if !ok {
			key = v.Type().Field(i).Name
		}
		switch val.Interface().(type) {
		case string:
			if !val.IsZero() {
				updates = append(updates, bson.E{Key: key, Value: val.String()})
			}
		case uint, uint8, int:
			if !val.IsZero() {
				updates = append(updates, bson.E{Key: key, Value: val.Int()})
			}
		case bool:
			updates = append(updates, bson.E{Key: key, Value: val.Bool()})
		}
	}
	return updates
}
