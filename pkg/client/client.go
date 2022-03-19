package client

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/dimitrijed93/dgtc/internal/handshake"
	"github.com/dimitrijed93/dgtc/internal/message"
	"github.com/dimitrijed93/dgtc/internal/peer"
	"github.com/dimitrijed93/dgtc/internal/utils"
)

// Wrapper for TCP connection with Peer
type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield utils.Bitfield
	peer     peer.Peer
	infoHash [utils.INFO_HASH_LEN]byte
	peerID   [utils.INFO_HASH_LEN]byte
}

func New(peer peer.Peer, peerId, infohash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout(utils.PROTOCOL, peer.String(), utils.CONNECTION_DEADLINE)
	if err != nil {
		return nil, err
	}

	err = completeHandshake(conn, infohash, peerId)

	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := utils.ReceiveBitField(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infohash,
		peerID:   peerId,
	}, nil

}

func completeHandshake(conn net.Conn, infohash, peerId [20]byte) error {
	conn.SetDeadline(time.Now().Add(utils.CONNECTION_DEADLINE))
	defer conn.SetDeadline(time.Time{})

	req := handshake.NewHandshake(infohash, peerId)

	_, err := conn.Write(req.Serialize())

	if err != nil {
		return err
	}

	res, err := handshake.Read(conn)

	if err != nil {
		return err
	}
	if !bytes.Equal(req.InfoHash[:], infohash[:]) {
		return fmt.Errorf("Expected infohash %x but got %x", res.InfoHash, infohash)
	}

	return nil

}

func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	return msg, err
}

func (c *Client) SendRequest(index, begin, length int) error {
	req := message.NewMessageRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())

	return err
}

func (c *Client) SendInterested() error {
	req := message.NewMessageInterested()
	_, err := c.Conn.Write(req.Serialize())
	return err
}

func (c *Client) SendNotInterested() error {
	req := message.NewMessageNotInterested()
	_, err := c.Conn.Write(req.Serialize())
	return err
}

func (c *Client) SendUnchoke() error {
	req := message.NewMessageUnchoke()
	_, err := c.Conn.Write(req.Serialize())
	return err
}

func (c *Client) SendHave(index int) error {
	req := message.NewMessageHave(index)
	_, err := c.Conn.Write(req.Serialize())
	return err
}
