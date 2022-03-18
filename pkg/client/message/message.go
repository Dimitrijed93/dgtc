package message

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type messageID uint8

const (
	MESSAGE_CHOKE          messageID = 0
	MESSAGE_UNCHOKE        messageID = 1
	MESSAGE_INTERESTED     messageID = 2
	MESSAGE_NOT_INTERESTED messageID = 3
	MESSAGE_HAVE           messageID = 4
	MESSAGE_BITFIELD       messageID = 5
	MESSAGE_REQUEST        messageID = 6
	MESSAGE_PIECE          messageID = 7
	MESSAGE_CANCEL         messageID = 8
)

type Message struct {
	MessageID messageID
	Payload   []byte
}

// Create message to request a block (part of the piece) from a peer
// To request a block a client needs to send MESSAGE_INTERESTED
// Message and to be UNCHOKED by the peer
func NewMessageRequest(index int, begin int, len int) *Message {
	payload := make([]byte, 12)

	// Index of the piece
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	// Offset of the block inside a piece
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	// Block length
	binary.BigEndian.PutUint32(payload[8:12], uint32(len))

	return &Message{MessageID: MESSAGE_REQUEST, Payload: payload}
}

func NewMessageInterested() *Message {
	return &Message{
		MessageID: MESSAGE_INTERESTED,
	}
}

func NewMessageNotInterested() *Message {
	return &Message{
		MessageID: MESSAGE_NOT_INTERESTED,
	}
}

func NewMessageUnchoke() *Message {
	return &Message{
		MessageID: MESSAGE_UNCHOKE,
	}
}

// Create a message MESSAGE_HAVE to declare to the peers
// that client has a block of data
func NewMessageHave(index int) *Message {
	payload := make([]byte, 4)

	binary.BigEndian.PutUint32(payload, uint32(index))

	return &Message{MessageID: MESSAGE_HAVE, Payload: payload}
}

func ParseMessageHave(msg *Message) (int, error) {

	if msg.MessageID != MESSAGE_HAVE {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", MESSAGE_HAVE, msg.MessageID)
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}

	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

func ParseMessagePiece(index int, buf []byte, msg *Message) (int, error) {

	if msg.MessageID != MESSAGE_PIECE {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", MESSAGE_PIECE, msg.MessageID)
	}

	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}

	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))

	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}

	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))

	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}

	data := msg.Payload[8:]

	if begin+len(data) >= len(buf) {
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}

	copy(buf[begin:], data)

	return len(data), nil

}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.MessageID)
	copy(buf[5:], m.Payload)

	return buf

}

func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)

	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)

	_, err = io.ReadFull(r, messageBuf)

	if err != nil {
		return nil, err
	}

	m := Message{
		MessageID: messageID(messageBuf[0]),
		Payload:   messageBuf[1:],
	}

	return &m, nil
}

func (m *Message) String() string {
	typeString := reflect.TypeOf(m).String()

	if m == nil {
		return typeString
	}

	return fmt.Sprintf("%s [%d]", typeString, len(m.Payload))
}
