package tracker

import (
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"net/url"
	"time"

	"github.com/dimitrijed93/dgtc/internal/files"
	"github.com/dimitrijed93/dgtc/internal/peer"
	"github.com/dimitrijed93/dgtc/internal/utils"
)

const (
	LAST_INDEX_OF_PROTOCOL                = 6
	TIMEOUT                               = time.Second * 2
	ERROR_FAIL_TO_CONSTRUCT_URL           = "Fail to construct UDP Url"
	ERROR_FAIL_TO_DIAL_URL                = "Fail to dial Url: "
	ERROR_FAIL_TO_WRITE_CONNECT           = "Fail to Write: "
	ERROR_INVALID_LEN                     = "Invalid len "
	ERROR_FAIL_TO_READ_CONNECT            = "Fail to set connection"
	ERROR_FAIL_TO_SEND_ANNOUNCE           = "Fail to send announce"
	ERROR_FAIL_TO_READ_ANNOUNCE           = "Fail to Read Announce Response"
	ERROR_CONNECT_INVALID_TRANSACTION_ID  = "Invalid Connect TransactionId"
	ERROR_ANNOUNCE_INVALID_TRANSACTION_ID = "Invalid Announce TransactionId"
	ERROR_INVALID_ANNOUNCE_ACTION         = "Invalid Announce Action"
	ERROR_INVALID_ANNOUNCE_ACTION_PEERS   = "Invalid Announce Action Peers"
	ERROR_INVALID_PEERS                   = "Invalid Peers"
	ERROR_INVALID_CONNECT_ACTION          = "Invalid Connect Action"
	ACTION_CONNECT                        = 0
	ACTION_ANNOUNCE                       = 1
)

type UdpTracker struct {
	Tf            files.TorrentFile
	conn          net.Conn
	connectionId  uint64
	transactionId uint32
}

func NewUdpTracker(tf files.TorrentFile) *UdpTracker {
	return &UdpTracker{
		Tf:            tf,
		transactionId: rand.Uint32(),
	}
}

func (tracker *UdpTracker) buildTrackerUrl(peerID [20]byte, port uint16) (*url.URL, error) {
	base, err := url.Parse(tracker.Tf.Announce)
	if err != nil {
		return &url.URL{}, err
	}
	return base, nil
}

func (tracker *UdpTracker) RequestPeers(peerId [20]byte) ([]peer.Peer, error) {
	url, err := tracker.buildTrackerUrl(peerId, 0)

	if err != nil {
		log.Fatal(ERROR_FAIL_TO_CONSTRUCT_URL)
		return nil, err
	}

	conn, err := net.DialTimeout(utils.PROTOCOL_UDP, url.Host, TIMEOUT)

	defer conn.Close()

	if err != nil {
		log.Print(ERROR_FAIL_TO_DIAL_URL, url)
		return nil, err
	}

	buf := tracker.createConnectPacket()

	n, err := conn.Write(buf)

	if err != nil {
		log.Print(ERROR_FAIL_TO_WRITE_CONNECT, buf)
		return nil, err
	}

	if n != len(buf) {
		log.Print(ERROR_INVALID_LEN, n)
		return nil, err
	}

	tracker.conn = conn.(*net.UDPConn)
	tracker.obtainConnectionId()
	tracker.announce()
	peers, err := tracker.obtainPeers()

	if err != nil {
		log.Print(ERROR_INVALID_PEERS)
	}
	return peers, nil
}

func (tracker *UdpTracker) announce() {
	req := make([]byte, 98)

	binary.BigEndian.PutUint64(req[0:], tracker.connectionId)
	binary.BigEndian.PutUint32(req[8:], 1)
	binary.BigEndian.PutUint32(req[12:], tracker.transactionId)

	peerId := utils.NewPeerId()

	copy(req[36:], peerId[:])
	copy(req[16:], tracker.Tf.InfoHash[:])
	binary.BigEndian.PutUint64(req[56:], 0)
	binary.BigEndian.PutUint64(req[64:], uint64(tracker.Tf.Length))
	binary.BigEndian.PutUint64(req[72:], 0)
	binary.BigEndian.PutUint32(req[80:], 0)
	binary.BigEndian.PutUint32(req[84:], 0)
	binary.BigEndian.PutUint32(req[88:], 21234)
	binary.BigEndian.PutUint32(req[92:], 1)
	binary.BigEndian.PutUint16(req[96:], utils.PORT)

	tracker.conn.SetWriteDeadline(time.Now().Add(TIMEOUT))

	_, err := tracker.conn.Write(req)

	if err != nil {
		log.Print(ERROR_FAIL_TO_SEND_ANNOUNCE, err)
	}
}

func (tracker *UdpTracker) obtainPeers() ([]peer.Peer, error) {
	bufRes := make([]byte, 1000)

	tracker.conn.SetReadDeadline(time.Now().Add(TIMEOUT))
	n, err := tracker.conn.Read(bufRes)
	if err != nil {
		log.Print(ERROR_FAIL_TO_READ_ANNOUNCE, err, n)
	}

	action := binary.BigEndian.Uint32(bufRes[0:4])
	transactionId := binary.BigEndian.Uint32(bufRes[4:8])
	leechers := binary.BigEndian.Uint32(bufRes[12:16])
	seeders := binary.BigEndian.Uint32(bufRes[16:20])
	addr := bufRes[20:]

	if action != ACTION_ANNOUNCE {
		log.Fatal(ERROR_INVALID_ANNOUNCE_ACTION_PEERS)
	}

	if transactionId != tracker.transactionId {
		log.Fatal(ERROR_ANNOUNCE_INVALID_TRANSACTION_ID)
	}

	addrLen := net.IPv4len
	step := addrLen + 2

	var peers []peer.Peer

	for i := step; i < len(addr); i += step {

		ip := make(net.IP, addrLen)
		copy(ip, addr[i-step:i-2])
		port := binary.BigEndian.Uint16(addr[i-2 : i])

		peer := peer.Peer{
			Port: port,
			IP:   ip,
		}
		peers = append(peers, peer)
		break

	}
	log.Print("Seeders ", seeders)
	log.Print("leechers ", leechers)

	return peers, nil

}

func (tracker *UdpTracker) obtainConnectionId() {

	bufRes := make([]byte, 16)

	n, err := tracker.conn.Read(bufRes)

	if err != nil {
		log.Print(ERROR_FAIL_TO_READ_CONNECT)
	}

	if n != 16 {
		log.Print(ERROR_INVALID_LEN, n)
	}

	action := binary.BigEndian.Uint32(bufRes[0:4])
	transactionId := binary.BigEndian.Uint32(bufRes[4:8])
	connectionId := binary.BigEndian.Uint64(bufRes[8:16])

	if action != ACTION_CONNECT {
		log.Fatal(ERROR_INVALID_CONNECT_ACTION)
	}

	if transactionId != tracker.transactionId {
		log.Fatal(ERROR_CONNECT_INVALID_TRANSACTION_ID)
	}

	tracker.connectionId = connectionId
}

func (tracker *UdpTracker) createConnectPacket() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[0:], utils.PROTOCOL_ID)      // magic constant
	binary.BigEndian.PutUint32(buf[8:], ACTION_CONNECT)         // action connect
	binary.BigEndian.PutUint32(buf[12:], tracker.transactionId) // transaction id

	return buf
}

func (tracker *UdpTracker) createAnnouncePacket() []byte {
	req := make([]byte, 98)

	binary.BigEndian.PutUint64(req[0:], tracker.connectionId)
	binary.BigEndian.PutUint32(req[8:], 1)
	binary.BigEndian.PutUint32(req[12:], tracker.transactionId)

	peerId := utils.NewPeerId()

	copy(req[36:], peerId[:])
	copy(req[16:], tracker.Tf.InfoHash[:])
	binary.BigEndian.PutUint64(req[56:], 0)
	binary.BigEndian.PutUint64(req[64:], uint64(tracker.Tf.Length))
	binary.BigEndian.PutUint64(req[72:], 0)
	binary.BigEndian.PutUint32(req[80:], 0)
	binary.BigEndian.PutUint32(req[84:], 0)
	binary.BigEndian.PutUint32(req[88:], 21234)
	binary.BigEndian.PutUint32(req[92:], 1)
	binary.BigEndian.PutUint16(req[96:], utils.PORT)

	return req
}
