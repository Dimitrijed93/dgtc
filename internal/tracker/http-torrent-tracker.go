package tracker

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dimitrijed93/dgtc/internal/files"
	"github.com/dimitrijed93/dgtc/internal/peer"
	"github.com/dimitrijed93/dgtc/internal/utils"
	"github.com/jackpal/bencode-go"
)

type HttpTracker struct {
	Tf files.TorrentFile
}

func (tracker *HttpTracker) buildTrackerUrl(peerID [20]byte, port uint16) (*url.URL, error) {
	base, err := url.Parse(tracker.Tf.Announce)
	if err != nil {
		return &url.URL{}, err
	}

	params := url.Values{
		// Identifies the file we’re trying to download.
		// It is calculated using bencode
		"info_hash": []string{string(tracker.Tf.InfoHash[:])},
		// Name of the client. Used to be identified
		// by other peers
		"port":       []string{strconv.Itoa(int(port))},
		"peer_id":    []string{string(peerID[:])},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tracker.Tf.Length)},
	}

	base.RawQuery = params.Encode()
	return base, nil
}

func (tracker *HttpTracker) RequestPeers(peerId [20]byte) ([]peer.Peer, error) {

	url, err := tracker.buildTrackerUrl(peerId, utils.PORT)

	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}

	res, err := c.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	trackerRes := TrackerResponse{}

	err = bencode.Unmarshal(res.Body, &trackerRes)
	if err != nil {
		return nil, err
	}

	return peer.Unmarshal([]byte(trackerRes.Peers))
}
