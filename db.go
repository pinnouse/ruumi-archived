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
CREATE TABLE IF NOT EXISTS nyaa_torrents (
	nid integer PRIMARY KEY,
	name varchar(255),
	torrent bytea
)`)
	return err
}

func addTorrent(db *sql.DB, id int32, name string, torrent []byte) (newId int32, err error) {
	statement := `
INSERT INTO nyaa_torrents (id, name, torrent)
VALUES ($1, $2, $3) RETURNING id`
	newId = 0
	err = db.QueryRow(statement, id, name, torrent).Scan(&newId)
	return newId, err
}

func getTorrentNyaa(db *sql.DB, id int32, name string, torrentData chan []byte) {
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
