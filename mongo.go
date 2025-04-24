package main

import (
	"context"
	"fmt"
	"log"

	// "go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoInstace struct {
	connectionInstance *mongo.Client
	collection *mongo.Collection
	database *mongo.Database
}

func createInstance(url string) MongoInstace {
	uri := "mongodb://localhost:27017"
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return MongoInstace{
		connectionInstance: client,
	}
}

func (m *MongoInstace) closeInstance(){
	if err := m.connectionInstance.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (m *MongoInstace) connectDB(db string){
	m.database = m.connectionInstance.Database(db)
}

func (m *MongoInstace) connectCollection(collection string){
	m.collection = m.database.Collection(collection)
}

func (m *MongoInstace) insertData(pageInfos []*PageInfo){
	docs := make([]interface{}, len(pageInfos))
	for i, p := range pageInfos {
		docs[i] = p
	}

	insertResult, err := m.collection.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted documents: %v\n", insertResult.InsertedIDs)
}