package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass"
	dbname   = "ruumi"
)

func connectServer() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
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

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS torrents (
	nid integer PRIMARY KEY,
	name varchar(255),
	torrent text
)`)
	return err
}

func addTorrent(db *sql.DB, id int32, name string, torrent string) (newId int32, err error) {
	statement := `
INSERT INTO torrents (id, name, torrent)
VALUES ($1, $2, $3) RETURNING id`
	newId = 0
	err = db.QueryRow(statement, id, name, torrent).Scan(&newId)
	return newId, err
}

func getTorrentNyaa(db *sql.DB, id int32, torrentData chan string) {
	statement := `SELECT torrent FROM torrents WHERE nid = $1 LIMIT 1`
	row := db.QueryRow(statement, id)
	tData := ""
	err := row.Scan(&tData)
	switch err {
	case sql.ErrNoRows:
		//Query for one
		fmt.Println("Looking for data")
		tData := request(fmt.Sprintf("https://nyaa.si/download/%d.torrent", id))
		//decode(tData)
		fmt.Printf("Got something! %d\n", len(tData))
		torrentData <- string(tData)
	default:
		fmt.Println(err)
		fmt.Println("How did we get here")
		torrentData <- tData
	}
}
