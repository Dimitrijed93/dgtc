package dgtc

import (
	"crypto/rand"
	"log"

	"github.com/dimitrijed93/dgtc/internal/downloader"
	"github.com/dimitrijed93/dgtc/internal/files"
	"github.com/dimitrijed93/dgtc/internal/tracker"
	"github.com/dimitrijed93/dgtc/internal/utils"
)

type Dgtc struct {
	In     string
	Out    string
	Type   string
	PeerId [20]byte
}

func NewDgtc(in string, out string) *Dgtc {
	var peerId [utils.PEER_ID_LEN]byte
	_, err := rand.Read(peerId[:])
	if err != nil {
		panic("Error generating peerdId")
	}

	return &Dgtc{
		In:     in,
		Out:    out,
		PeerId: peerId,
	}
}

func (d *Dgtc) Start() {
	tf, err := files.NewTorrentFile(d.In)

	tracker := d.determineTracker(tf)

	peers, err := tracker.RequestPeers(d.PeerId)

	if err != nil {
		log.Fatal(err)
	}

	dw, err := downloader.NewDownloader(tf, peers, d.PeerId)

	dw.Start(d.Out)
	if err != nil {
		log.Fatal(err)
	}
}

func (d Dgtc) determineTracker(tf files.TorrentFile) tracker.Tracker {
	if tf.IsHTTP {
		return &tracker.HttpTracker{
			Tf: tf,
		}
	} else {
		return tracker.NewUdpTracker(
			tf,
		)
	}
}
