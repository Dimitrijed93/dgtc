package tracker

import (
	"github.com/dimitrijed93/dgtc/internal/files"
	"github.com/dimitrijed93/dgtc/internal/peer"
	"github.com/dimitrijed93/dgtc/internal/utils"
)

type UdpTracker struct {
	tf files.TorrentFile
}

func (tracker *UdpTracker) buildTrackerUrl(peerID [20]byte, port uint16) (string, error) {
	// NOOP
	return utils.EMPTY_STRING, nil
}

func (tracker *UdpTracker) RequestPeers(peerId [20]byte) ([]peer.Peer, error) {
	// NOOP
	return nil, nil
}
