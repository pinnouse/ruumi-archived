package main

import (
	"fmt"
	bencode "github.com/IncSW/go-bencode"
	"github.com/mitchellh/mapstructure"
)

func decode(torrent []byte) (data map[string]interface{}) {
	i, err := bencode.Unmarshal(torrent)
	if err != nil {
		panic(err)
	}

	data = make(map[string]interface{})
	switch v := i.(type) {
	case map[string]interface{}:
		for key, val := range v {
			data[key] = val
		}
	}
	return data
}

func decodeBT(bytes []byte) (torrent BitTorrent, err error) {
	i, err := bencode.Unmarshal(bytes)
	fmt.Println(i)
	torrent = BitTorrent{}
	mapstructure.Decode(i, &torrent)
	return
}

func getFiles(torrent BitTorrent) {

}
