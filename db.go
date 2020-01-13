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

func dbSearchCollection(client *mongo.Client) (collection *mongo.Collection) {
	return client.Database("gogo").Collection("searches")
}

func dbCategoryCollection(client *mongo.Client) (collection *mongo.Collection) {
	return client.Database("gogo").Collection("categories")
}

func dbEpisodeCollection(client *mongo.Client) (collection *mongo.Collection) {
	return client.Database("gogo").Collection("episodes")
}

func dbSetSearch(client *mongo.Client, results GOGOSearchResults) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := dbSearchCollection(client)
	r := GOGOSearchResults{}
	filter := bson.M{"searchTerm": results.SearchTerm}
	err := collection.FindOne(ctx, filter).Decode(&r)
	if err == nil {
		_, err = collection.ReplaceOne(ctx, filter, bson.M{
			"searchTerm":  results.SearchTerm,
			"results":     results.Results,
			"lastUpdated": time.Now().Unix(),
		})
	} else {
		_, err = collection.InsertOne(ctx, bson.M{
			"searchTerm":  results.SearchTerm,
			"results":     results.Results,
			"lastUpdated": time.Now().Unix(),
		})
	}
	fmt.Println(err)
	return err == nil
}

func dbGetSearch(client *mongo.Client, searchTerm string, page int) (results GOGOSearchResults) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := dbSearchCollection(client)
	filter := bson.M{
		"searchTerm": searchTerm,
	}
	err := collection.FindOne(ctx, filter).Decode(&results)
	if err != nil || time.Now().Unix()-results.LastUpdated > 60*60*24 {
		categories := make(chan []GOGOCategory)
		go gogoSearch(searchTerm, page, categories)
		results = GOGOSearchResults{
			SearchTerm:  searchTerm,
			Results:     <-categories,
			LastUpdated: 0,
		}
		dbSetSearch(client, results)
	}
	return
}

func dbSetCategory(client *mongo.Client, category GOGOCategory) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := dbCategoryCollection(client)
	updated := bson.M{
		"name":     category.Name,
		"catURL":   category.CatURL,
		"episodes": category.Episodes,
	}
	cat, err := dbGetCategory(client, category.CatURL)
	if err == nil && cat.Episodes < category.Episodes {
		filter := bson.M{"catURL": category.CatURL}
		_, err = collection.ReplaceOne(ctx, filter, updated)
	} else if err != nil {
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
	collection := client.Database("gogo").Collection("episodes")
	_, err := collection.InsertOne(
		ctx, bson.M{
			"epNum":       episode.EpNum,
			"srcURL":      episode.SrcURL,
			"category":    episode.Category,
			"lastUpdated": time.Now().Unix(),
		})
	return err == nil
}

func dbGetEpisode(client *mongo.Client, categoryURL string, episodeNum int) (episode GOGOEpisode, err error) {
	filter := bson.M{"category": categoryURL, "epNum": episodeNum}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = dbEpisodeCollection(client).FindOne(ctx, filter).Decode(&episode)
	if err != nil || time.Now().Unix()-episode.LastUpdated > 60*5 {
		epSrcChan := make(chan string)
		go gogoFetchEpisode(categoryURL, episodeNum, epSrcChan)
		epSrc := <-epSrcChan
		episode = GOGOEpisode{
			EpNum:       episodeNum,
			SrcURL:      epSrc,
			Category:    categoryURL,
			LastUpdated: time.Now().Unix(),
		}
		dbSetEpisode(client, episode)
	}
	return episode, nil
}
