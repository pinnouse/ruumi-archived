package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	dbHost = "mongodb://localhost:27017/?readPreference=primary&ssl=false"
)

func connectDB() (client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbHost))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return
}

func dbCollection(client *mongo.Client, collectionName string) (collection *mongo.Collection) {
	return client.Database("ruumi").Collection(collectionName)
}

func getUser(client *mongo.Client, userId string) (user User, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = dbCollection(client, "users").FindOne(ctx, bson.M{"id": userId}).Decode(&user)
	return
}

func addUser(client *mongo.Client, user User) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = dbCollection(client, "users").InsertOne(ctx, user)
	return
}

func search(client *mongo.Client, query string) (results []Anime, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := dbCollection(client, "anime").Find(ctx, bson.D{
		{"$or", bson.A{
			bson.D{
				{"title", primitive.Regex{Pattern: query, Options: "i"}},
			},
			bson.D{
				{"altTitles", primitive.Regex{Pattern: query, Options: "i"}},
			},
		}},
	})
	if err != nil {
		log.Println(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result Anime
		if err = cur.Decode(&result); err != nil {
			log.Println(err)
		}
		results = append(results, result)
	}
	if err = cur.Err(); err != nil {
		log.Println(err)
	}
	return
}

func getAnime(client *mongo.Client, id int32) (result Anime, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = dbCollection(client, "anime").FindOne(
		ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	return
}

func addAnime(client *mongo.Client, anime Anime) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = dbCollection(client, "anime").InsertOne(ctx, anime)
	return
}

func addEpisode(client *mongo.Client, animeId int32, episode Episode) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = dbCollection(client, "anime").UpdateOne(
		ctx,
		bson.M{"id": animeId},
		bson.D{{
			"episodes",
			bson.D{{"$addToSet", episode}}},
		})
	return
}
