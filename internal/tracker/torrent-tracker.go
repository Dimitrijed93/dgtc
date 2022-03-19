package tracker

import (
	"github.com/dimitrijed93/dgtc/internal/peer"
)

const (
	HTTP = "http"
	UDP  = "udp"
)

type Tracker interface {
	buildTrackerUrl(peerID [20]byte, port uint16) (string, error)
	RequestPeers(peerId [20]byte) ([]peer.Peer, error)
}

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}
