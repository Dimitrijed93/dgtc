package client

import (
	"fmt"
	"net"
	"time"

	"github.com/dimitrijed93/dgtc/pkg/bitfield"
	"github.com/dimitrijed93/dgtc/pkg/message"
	"github.com/dimitrijed93/dgtc/pkg/peer"
	"github.com/dimitrijed93/dgtc/pkg/utils"
)

/** Wrapper for TCP connection with Peer */
type Connection struct {
	Conn     net.Conn
	Choked   bool
	Bitfield bitfield.Bitfield
	peer     peer.Peer
	infoHash [utils.INFO_HASH_LEN]byte
	peerID   [utils.INFO_HASH_LEN]byte
}

func receiveBitField(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(utils.CONNECTION_DEADLINE)) // Add Timeout
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)

	if err != nil {
		return nil, err
	}

	// Accept only bitfield messages
	if msg.MessageID != message.MESSAGE_BITFIELD {
		err := fmt.Errorf("Expected bitfield but got ID %d", msg.MessageID)
		return nil, err
	}

	return msg.Payload, nil
}

/** Setup new Connection with peer and completes Handshake */
// func NewConnection(peer peer.Peer, peerID, infoHash [utils.INFO_HASH_LEN]byte) (*Connection, error) {

// }
