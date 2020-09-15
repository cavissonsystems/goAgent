package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	nd "goAgent"
	logger "goAgent/logger"
	"goAgent/module/cavmongo"
	"log"
	"time"
)

func m3(bt uint64) {
	nd.Method_entry(bt, "a.b.m3")
	time.Sleep(2 * time.Millisecond)
	nd.Method_exit(bt, "a.b.m3")
}

func m2(bt uint64) {
	nd.Method_entry(bt, "a.b.m2")
	time.Sleep(2 * time.Millisecond)
	nd.Method_exit(bt, "a.b.m2")
}

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func Call_mongo(ctx context.Context) {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(cavmongo.CommandMonitor())
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		logger.ErrorPrint("error come from client")
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {

		logger.ErrorPrint("Error : come from Ping")
	}

	logger.TracePrint("Connected to MongoDB!")

	collection := client.Database("test").Collection("trainers")

	ash := Trainer{"Ash", 10, "Pallet Town"}

	insertResult, err := collection.InsertOne(ctx, ash)
	if err != nil {
		log.Fatal(err)
	}

	//logger.TracePrint("Inserted a single document: ", insertResult.InsertedID)
	log.Println("Inserted a single document: ", insertResult.InsertedID)
	filter := bson.D{{"name", "Ash"}}

	var result Trainer

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	//logger.TracePrint("Found a single document:%+v ", result)
	log.Println("Found a single document:%+v ", result)

	bt := ctx.Value("CavissonTx").(uint64)

	m2(bt)

	m3(bt)
}
