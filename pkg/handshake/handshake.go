package handshake

import (
	"fmt"
	"io"

	"github.com/dimitrijed93/dgtc/pkg/utils"
)

type Handshake struct {
	Pstr     string   // Type of Protocol
	InfoHash [20]byte // Hash of the file
	PeerID   [20]byte // Peer identity
}

func (h *Handshake) NewHandshake(infoHash [20]byte, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

/** Serialize Handshake constisting of :
Length of the protocol
Protocol
Optional Bitfields (8 bytes)
InfoHash of the file we want to download
PeerId
*/
func (h *Handshake) Serialize() []byte {
	pstrlen := len(h.Pstr)
	bufLen := 49 + pstrlen
	buf := utils.NewBuf(bufLen)
	buf[0] = byte(pstrlen)
	copy(buf[1:], h.Pstr)
	// Leave 8 reserved bytes
	copy(buf[1+pstrlen+8:], h.InfoHash[:])
	copy(buf[1+pstrlen+8+20:], h.PeerID[:])
	return buf

}

func Read(r io.Reader) (*Handshake, error) {
	lenBuf := utils.NewBuf(1)

	_, err := io.ReadFull(r, lenBuf)

	if err != nil {
		return nil, err
	}

	pstrLen := int(lenBuf[0])

	if pstrLen == 0 {
		err = fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	handshakeBuf := utils.NewBuf(48 + pstrLen)

	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash [20]byte
	var peerID [20]byte

	// Add bitfield lenght to Pstr
	infoHashStart := pstrLen + utils.BITFIELD_LEN
	peerIDStart := infoHashStart + utils.INFO_HASH_LEN

	copy(infoHash[:], handshakeBuf[infoHashStart:peerIDStart])
	copy(peerID[:], handshakeBuf[peerIDStart:])

	hanshake := Handshake{
		Pstr:     string(handshakeBuf[0:pstrLen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &hanshake, nil
}
