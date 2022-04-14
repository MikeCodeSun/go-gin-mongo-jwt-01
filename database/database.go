package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func DbInstance() *mongo.Client {
	// get .env varible
	if err := godotenv.Load(); err != nil {
		log.Panic(err.Error())
	}
	uri := os.Getenv("URI")
  // context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
	// mongo connect
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err !=nil {
		log.Panic(err)
	}
	defer cancel()

  return c
}

var Client = DbInstance()

func ConnectCollection(client *mongo.Client,name string) *mongo.Collection{
	collection := client.Database("go-jwt-auth").Collection(name)
	return collection
}

