package main

import (
	"net/url"
	"strconv"

	"github.com/dimitrijed93/dgtc/pkg/utils"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (t *TorrentFile) buildTrackerUrl(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return utils.EMPTY_STRING, err
	}

	params := url.Values{
		// Identifies the file we’re trying to download.
		// It is calculated using bencode
		"info_hash": []string{string(t.InfoHash[:])},
		// Name of the client. Used to be identified
		// by other peers
		"peer_id":    []string{string(peerID[:])},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}
