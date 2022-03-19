package dgtc

import (
	"crypto/rand"
	"log"

	"github.com/dimitrijed93/dgtc/internal/downloader"
	"github.com/dimitrijed93/dgtc/internal/files"
)

type Dgtc struct {
	In     string
	Out    string
	PeerId [20]byte
}

func NewDgtc(in string, out string) *Dgtc {
	var peerId [20]byte
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
	peers, err := tf.RequestPeers(d.PeerId)

	if err != nil {
		log.Fatal(err)
	}

	dw, err := downloader.NewDownloader(tf, peers, d.PeerId)

	dw.Init(d.In)
	if err != nil {
		log.Fatal(err)
	}
}
