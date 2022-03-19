package files

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/dimitrijed93/dgtc/internal/utils"
	"github.com/jackpal/bencode-go"
)

type TorrentInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

func (b *TorrentInfo) hash() ([20]byte, error) {
	var buff bytes.Buffer

	err := bencode.Marshal(&buff, *b)
	if err != nil {
		return [20]byte{}, err
	}

	hash := sha1.Sum(buff.Bytes())
	return hash, nil
}

func (b *TorrentInfo) splitHashes() ([][20]byte, error) {
	buf := []byte(b.Pieces)
	if len(buf)%utils.INFO_HASH_LEN != 0 {
		return nil, fmt.Errorf("Malformed pieces of lenght %d", len(buf))
	}

	numHashes := len(buf) / utils.INFO_HASH_LEN
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*utils.INFO_HASH_LEN:(i+1)*utils.INFO_HASH_LEN])
	}

	return hashes, nil
}
