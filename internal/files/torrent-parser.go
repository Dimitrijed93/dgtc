package files

import (
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type ParsedTorrentFile struct {
	Announce string      `bencode:"announce"`
	Info     TorrentInfo `bencode:"info"`
}

// Parses Meta-Info file from a file system into ParsedTorrentFile
func ParseMetaInfo(path string) (ParsedTorrentFile, error) {
	file, err := os.Open(path)
	fmt.Print("file", file)
	if err != nil {
		return ParsedTorrentFile{}, err
	}
	defer file.Close()

	bto := ParsedTorrentFile{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return ParsedTorrentFile{}, err
	}

	return bto, nil
}
