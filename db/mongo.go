package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"time"
)

type ClientParams struct {
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
}

func ConnectMongoDatabase(params ClientParams) (*mongo.Database, error) {
	var db *mongo.Database

	dbAuthCreds := options.Credential{
		AuthMechanism: auth.SCRAMSHA256,
		Username:      params.DBUser,
		Password:      params.DBPassword,
	}
	dbClientOptions := options.Client()

	dbClientOptions.SetDirect(true)
	dbClientOptions.SetAuth(dbAuthCreds)
	dbClientOptions.ApplyURI(params.DBHost)

	dbClient, err := mongo.NewClient(dbClientOptions)

	if err != nil {
		return db, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()

	err = dbClient.Connect(ctx)

	if err != nil {
		return db, err
	}

	db = dbClient.Database(params.DBName)

	return db, nil
}
