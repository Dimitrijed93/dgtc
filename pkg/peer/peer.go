package peer

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

/** Container for Peer indentity
Usual Peer ports: 6881-6889 TCP
Peers are registered on the Tracker
*/
type Peer struct {
	IP   net.IP
	Port uint16
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func Unmarshal(peersBin []byte) ([]Peer, error) {
	// 2 bytes for port, 4 bytes for IP
	const peerSize = 6

	numPeers := len(peersBin) / peerSize

	if len(peersBin)%numPeers != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}

	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize

		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])

	}

	return peers, nil
}
