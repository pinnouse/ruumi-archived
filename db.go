package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func connectDB() (client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(
		fmt.Sprintf(
			"mongodb://%s:27017/?readPreference=primary&ssl=false",
			func() string {
				if len(os.Getenv("DB_HOST")) > 0 {
					return os.Getenv("DB_HOST")
				} else {
					return "localhost"
				}
			}(),
		)))
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
	if err = cur.All(ctx, &results); err != nil {
		log.Println(err)
	}
	if err = cur.Err(); err != nil {
		log.Println(err)
	}
	return
}

func getAnime(client *mongo.Client, id string) (result Anime, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return
	}
	err = dbCollection(client, "anime").FindOne(
		ctx, bson.M{"_id": objectId}).Decode(&result)
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

func addEpisode(client *mongo.Client, animeId string, episode Episode) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objectId, err := primitive.ObjectIDFromHex(animeId)
	if err != nil {
		return
	}
	_, err = dbCollection(client, "anime").UpdateOne(
		ctx,
		objectId,
		bson.D{{
			"episodes",
			bson.D{{"$addToSet", episode}}},
		})
	return
}
