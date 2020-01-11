package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	db_host = "mongodb://localhost:27017/?readPreference=primary&ssl=false"
)

func connectDB() (client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(db_host))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	return client
}

func dbCategoryCollection(client *mongo.Client) (collection *mongo.Collection) {
	return client.Database("gogo").Collection("categories")
}

func dbEpisodeCollection(client *mongo.Client) (collection *mongo.Collection) {
	return client.Database("gogo").Collection("episodes")
}

func dbSetCategory(client *mongo.Client, category GOGOCategory) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := dbCategoryCollection(client)
	updated := bson.M{
		"name":     category.Name,
		"catURL":   category.CatURL,
		"episodes": category.Episodes,
	}
	_, err := dbGetCategory(client, category.CatURL)
	if err == nil {
		filter := bson.M{"catURL": category.CatURL}
		_, err = collection.UpdateOne(ctx, filter, updated)
	} else {
		_, err = collection.InsertOne(ctx, updated)
	}
	return err == nil
}

func dbGetCategory(client *mongo.Client, categoryURL string) (category GOGOCategory, err error) {
	cat := GOGOCategory{}
	filter := bson.M{"catURL": categoryURL}
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	collection := dbCategoryCollection(client)
	err = collection.FindOne(ctx, filter).Decode(&cat)
	return cat, err
}

func dbSetEpisode(client *mongo.Client, episode GOGOEpisode) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := dbEpisodeCollection(client)
	_, err := collection.InsertOne(
		ctx, bson.M{
			"epNum":    episode.EpNum,
			"srcURL":   episode.SrcURL,
			"category": episode.Category,
		})
	fmt.Println(err)
	return err == nil
}

func dbGetEpisode(client *mongo.Client, categoryURL string, episodeNum int) (episode GOGOEpisode, err error) {
	filter := bson.M{"category": categoryURL, "epNum": episodeNum}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = dbEpisodeCollection(client).FindOne(ctx, filter).Decode(&episode)
	if err != nil {
		fmt.Println(err)
		epSrcChan := make(chan string)
		go gogoFetchEpisode(categoryURL, episodeNum, epSrcChan)
		epSrc := <-epSrcChan
		episode = GOGOEpisode{
			EpNum:    episodeNum,
			SrcURL:   epSrc,
			Category: categoryURL,
		}
		dbSetEpisode(client, episode)
	}
	return
}
