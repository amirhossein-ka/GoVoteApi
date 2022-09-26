package mongo

import (
	"GoVoteApi/config"
	"GoVoteApi/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoImpl struct {
	// client *mongo.Client
	db      *mongo.Database
	cfg     *config.DBMongo
	userCol *mongo.Collection
	voteCol *mongo.Collection
}

func New(ctx context.Context, cfg *config.DBMongo) (repository.Repository, error) {
	ctx, cancel := newMongoContext(ctx, 10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	database := client.Database(cfg.DBName)

	return &mongoImpl{
		// client: client,
		db:  database,
		cfg: cfg,
        userCol: database.Collection("users"),
        voteCol: database.Collection("votes"),
	}, nil
}

func newMongoContext(ctx context.Context, timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
}
