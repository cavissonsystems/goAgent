package main

import (
    "context"
   logger "goAgent/logger"
   "log"
    "goAgent/module/cavmongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

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
}
