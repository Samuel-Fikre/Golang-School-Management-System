package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	MongoDB := "mongodb://localhost:27107"
	fmt.Print(MongoDB)

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDB))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
}
