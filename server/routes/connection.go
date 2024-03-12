package routes

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func DBinstance() *mongo.Client {
// 	MongoDb := "mongodb://localhost:27017/caloriesdb"
// 	/* client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
// 	if err != nil {
// 		log.Fatal(err)
// 	} */

// 	clientOptions := options.Client().ApplyURI(MongoDb)

//     // Connect to MongoDB
//     client, err := mongo.Connect(context.TODO(), clientOptions)

//     if err != nil {
//         log.Fatal(err)
//     }

// 	fmt.Println("Conectado a mongodb")

// 	return client
// }

func DBinstance() *mongo.Client {
	MongoDb := "mongodb://localhost:27017/caloriesdb"
	clientOptions := options.Client().ApplyURI(MongoDb)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conectado a mongodb")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("caloriesdb").Collection(collectionName)
	return collection
}
