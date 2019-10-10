package repository

import (
	"context"
	"easymile-tracker/common/database"
	"easymile-tracker/common/models"
	"fmt"
	"log"
)

// Save the trip
func Save(trip models.Trip) error {
	client := database.Connection()
	collection := client.Database("easymile").Collection("trip")
	insertResult, err := collection.InsertOne(context.TODO(), trip)

	if err != nil {
		fmt.Print(err)
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")

	if err != nil {
		return err
	}

	return nil
}
