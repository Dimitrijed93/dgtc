package peer

import (
	"encoding/binary"
	"fmt"
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

	numPeers := len(peersBin) / utils.PEER_SIZE

	if len(peersBin)%numPeers != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}

	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * utils.PEER_SIZE

		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])

	}

	return peers, nil
}
