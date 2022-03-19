package peer

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"

	"github.com/dimitrijed93/dgtc/internal/utils"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func Unmarshal(peersBin []byte) ([]Peer, error) {

	numPeers := validatePeersBin(peersBin)

	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * utils.PEER_SIZE
		// 4 for host, 2 for port
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])

	}

	return peers, nil
}

func validatePeersBin(peersBin []byte) int {

	if utils.IsEmpty(peersBin) {
		log.Fatal("There are no peers present")
	}

	numPeers := len(peersBin) / utils.PEER_SIZE

	if len(peersBin)%numPeers != 0 {
		log.Fatal("Received malformed peers")
	}

	return numPeers
}
