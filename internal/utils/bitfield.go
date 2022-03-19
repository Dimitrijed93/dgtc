package utils

import (
	"fmt"
	"net"
	"time"

	"github.com/dimitrijed93/dgtc/internal/message"
)

// Bitfield indicates what pieces the peer has
//
// Present pieces are indicated by 1.
// Missing pieces are indicates by 0.
type Bitfield []byte

func (b Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	return b[byteIndex]>>(7-offset)&1 != 0
}

func (b Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	b[byteIndex] |= 1 << (7 - offset)
}

func ReceiveBitField(conn net.Conn) (Bitfield, error) {
	conn.SetDeadline(time.Now().Add(CONNECTION_DEADLINE))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)

	if err != nil {
		return nil, err
	}

	if msg.MessageID != message.MESSAGE_BITFIELD {
		err := fmt.Errorf("Expected bitfield but got ID %d", msg.MessageID)
		return nil, err
	}

	return msg.Payload, nil
}
