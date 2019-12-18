package main

import (
	bencode "github.com/IncSW/go-bencode"
)

func decode(torrent []byte) (data interface{}) {
	i, err := bencode.Unmarshal(torrent)
	if err != nil {
		panic(err)
	}
	return i
}
