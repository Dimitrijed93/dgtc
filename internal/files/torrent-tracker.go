package files

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dimitrijed93/dgtc/internal/peer"
	"github.com/dimitrijed93/dgtc/internal/utils"
	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type trackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func NewTorrentFile(path string) (TorrentFile, error) {

	bto, err := ParseMetaInfo(path)
	infohash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	pieceHashes, err := bto.Info.splitHashes()
	if err != nil {
		return TorrentFile{}, nil
	}

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infohash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}

	return t, nil

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
		"port":       []string{strconv.Itoa(int(port))},
		"peer_id":    []string{string(peerID[:])},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) RequestPeers(peerId [20]byte) ([]peer.Peer, error) {

	url, err := t.buildTrackerUrl(peerId, utils.PORT)

	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}

	res, err := c.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	trackerRes := trackerResponse{}

	err = bencode.Unmarshal(res.Body, &trackerRes)
	if err != nil {
		return nil, err
	}

	return peer.Unmarshal([]byte(trackerRes.Peers))
}
