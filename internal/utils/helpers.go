package utils

import "crypto/rand"

func NewBuf(len int) []byte {
	return make([]byte, len)
}

func IsEmpty(list []byte) bool {
	len := len(list)
	return len == 0
}

func NewPeerId() [20]byte {
	var peerId [PEER_ID_LEN]byte
	_, err := rand.Read(peerId[:])
	if err != nil {
		panic("Error generating peerdId")
	}
	return peerId
}
