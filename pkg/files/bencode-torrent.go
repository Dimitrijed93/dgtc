package files

import (
	"os"

	"github.com/jackpal/bencode-go"
)

type BencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     BencodeInfo `bencode:"info"`
}

// Parses Meta-Info file from a file system into BencodeTorrent
func parseMetaInfo(path string) (BencodeTorrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return BencodeTorrent{}, err
	}
	defer file.Close()

	bto := BencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return BencodeTorrent{}, err
	}

	return bto, nil
}
