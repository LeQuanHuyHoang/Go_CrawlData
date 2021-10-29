package database

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Create a singleton object of MongoDB client
var clientInstance *mongo.Client

//Used to check error during creation
var clientInstanceError error

//Used to execute client creation proceduce (only one)
var mongoOnce sync.Once

const (
	CONNECTIONSTRING = "mongodb://localhost:27017"
	DB               = "intern_db"
)

//Return mongodb connection
func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once
	mongoOnce.Do(func() {
		//set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		//connect to mongodb
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}

		//check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = mongo.ErrClientDisconnected
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}
