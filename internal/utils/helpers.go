package utils

func NewBuf(len int) []byte {
	return make([]byte, len)
}

func IsEmpty(list []byte) bool {
	len := len(list)
	return len == 0
}
