package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass"
	dbname   = "ruumi"
)

func connectServer() *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	err = createTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

func createTable(db *sqlx.DB) (err error) {
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS nyaa_torrents (
	nid integer PRIMARY KEY,
	name varchar(255),
	torrent bytea
)`)
	if err != nil {
		return
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS gogo_categories (
	gid SERIAL PRIMARY KEY,
	name varchar(225),
	caturl varchar(255),
	episodes json
)`)
	if err != nil {
		return
	}
	return
}

func addCategory(db *sqlx.DB, catName string, catURL string) {
	id := 0
	err := db.QueryRow(`
INSERT INTO gogo_categories(name, caturl, episodes)
VALUES ($1, $2, '{}'::json) RETURNING gid
`, catName, catURL).Scan(&id)
	if err != nil {
		panic(err)
	}
}

func getEpisode(db *sqlx.DB, catName string, episodeNumber int, episode chan GOGOEpisode) {
	category := GOGOCategoryD{}
	err := db.Get(&category, `SELECT * FROM gogo_categories WHERE name = $1 LIMIT 1`, catName)
	if err != nil {
		fmt.Println(err)
		return
	}

	eps := make(map[string]interface{})
	err = json.Unmarshal([]byte(category.Episodes), &eps)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ep, ok := eps[strconv.Itoa(episodeNumber)]; ok {
		episodeHolder := GOGOEpisode{}
		marshalizedEp, err := json.Marshal(ep)
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal(marshalizedEp, &episodeHolder)
		episode <- episodeHolder
	} else {
		epSrc := make(chan string)
		go gogoFetchEpisode(category.CatURL, episodeNumber, epSrc)
		src := <-epSrc
		episode <- GOGOEpisode{
			EpID:   strconv.Itoa(episodeNumber),
			SrcURL: src,
		}
	}
}

func addTorrent(db *sqlx.DB, id int32, name string, torrent []byte) (newId int32, err error) {
	statement := `
INSERT INTO nyaa_torrents (id, name, torrent)
VALUES ($1, $2, $3) RETURNING id`
	newId = 0
	err = db.QueryRow(statement, id, name, torrent).Scan(&newId)
	return newId, err
}

func getTorrentNyaa(db *sqlx.DB, id int32, name string, torrentData chan []byte) {
	statement := `SELECT torrent FROM nyaa_torrents WHERE nid = $1 LIMIT 1`
	row := db.QueryRow(statement, id)
	tData := []byte{}
	err := row.Scan(&tData)
	switch err {
	case sql.ErrNoRows:
		//Query for one
		fmt.Printf("Adding to DB %d\n", id)
		tData, err := request(fmt.Sprintf("https://nyaa.si/download/%d.torrent", id))
		if err != nil {
			//TODO: Handler for not being able to download the torrent
			panic(err)
		}
		//decode(tData)
		addTorrent(db, id, name, tData)
		torrentData <- tData
	default:
		fmt.Println(err)
		fmt.Println("How did we get here")
		torrentData <- tData
	}
}
