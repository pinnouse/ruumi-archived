package main

type BTFile struct {
	Length uint32
	Path   [][]byte
}

type BTInfo struct {
	Length      uint32
	Files       []BTFile
	Name        []byte
	PieceLength uint32
	Pieces      []byte
}

type BitTorrent struct {
	Announce     []byte     `mapstructure:"announce"`
	AnnounceList [][][]byte `mapstructure:"announce-list"`
	Comment      []byte     `mapstructure:"comment"`
	CreatedBy    []byte     `mapstructure:"created by"`
	CreationDate uint32     `mapstructure:"creation date"`
	Encoding     []byte     `mapstructure:"encoding"`
	Info         BTInfo     `mapstructure:"info"`
}

type NyaaTorrent struct {
	Id   string
	Name string
}

//ENDPOINT STRUCTS
type ENDsearch struct {
	Torrents []NyaaTorrent
	LastPage bool
}
