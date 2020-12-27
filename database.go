package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type Resident struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"-" bson:"name"`
	Room     string             `json:"room" bson:"room"`
	Password string             `json:"-" bson:"password"`
}

type Booking struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Time int                `json:"time"`
	User Resident           `json:"user"`
}

func getBookings() (bookings []Booking) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer client.Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := client.Database("fasanhuset").Collection("washing").Find(ctx, bson.M{})
	var booking Booking
	for cursor.TryNext(context.Background()) {
		cursor.Decode(&booking)
		bookings = append(bookings, booking)
	}
	defer client.Disconnect(ctx)
	return bookings
}

func addBooking(booking Booking) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer client.Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"time": booking.Time}
	cursor, err := client.Database("fasanhuset").Collection("washing").Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}

	ok := cursor.TryNext(context.Background())
	if ok {
		fmt.Println("Exist already")
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_, err = client.Database("fasanhuset").Collection("washing").InsertOne(ctx, booking)
	if err != nil {
		fmt.Println(err)
	}
}

func removeBooking(booking Booking) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer client.Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"_id": booking.Id}
	_, err = client.Database("fasanhuset").Collection("washing").DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
}

func updateBooking(booking Booking) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer client.Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"_id": booking.Id}
	update := bson.M{"$set": booking}
	_, err = client.Database("fasanhuset").Collection("washing").UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
}

func getUser(id string) (user *Resident, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Disconnect(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objectId,
	}

	cursor, err := client.Database("fasanhuset").Collection("users").Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	cursor.TryNext(context.Background())

	err = cursor.Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (self *Resident) authenticate() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	defer client.Disconnect(ctx)

	if err != nil {
		fmt.Println(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{
		"room":     self.Room,
		"password": self.Password,
	}

	cursor, err := client.Database("fasanhuset").Collection("users").Find(ctx, filter)

	if err != nil {
		return err
	}

	cursor.TryNext(context.Background())

	err = cursor.Decode(&self)

	if err != nil {
		return err
	}

	return nil
}
