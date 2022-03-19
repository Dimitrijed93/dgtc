package utils

import "time"

const (
	EMPTY_STRING        = ""
	PROTOCOL            = "tcp"
	BITFIELD_LEN        = 8
	INFO_HASH_LEN       = 20
	PEER_ID_LEN         = 20
	PORT                = 6881
	PEER_SIZE           = 6 // 4 for host, 2 for port
	CONNECTION_DEADLINE = time.Second * 5
	PIECE_READ_DEADLINE = time.Second * 30
	MAX_BACKLOG         = 5 // Max number of pending requests
	MAX_BACKLOG_SIZE    = 16384
)
