package downloader

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/dimitrijed93/dgtc/internal/files"
	"github.com/dimitrijed93/dgtc/internal/message"
	"github.com/dimitrijed93/dgtc/internal/peer"
	"github.com/dimitrijed93/dgtc/internal/utils"
	"github.com/dimitrijed93/dgtc/pkg/client"
)

type Piece struct {
	index  int
	hash   [20]byte
	length int
}

type PieceResult struct {
	index int
	buf   []byte
}

type PieceJob struct {
	index  int
	hash   [20]byte
	length int
}

type PieceProgress struct {
	index      int
	client     *client.Client
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

type Downloader struct {
	Peers  []peer.Peer
	PeerId [20]byte
	T      files.TorrentFile
}

func NewDownloader(t files.TorrentFile, peers []peer.Peer, peerId [20]byte) (*Downloader, error) {

	downloader := &Downloader{
		Peers:  peers,
		PeerId: peerId,
		T:      t,
	}

	return downloader, nil
}

func (d *Downloader) Start(path string) error {

	buf, err := d.Download()
	if err != nil {
		return err
	}

	writeToOutFile(path, buf)
	return nil
}

func writeToOutFile(path string, buf []byte) {
	outFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Downloader) setupJobs(queue chan *PieceJob, results chan *PieceResult) {

	for index, hash := range d.T.PieceHashes {
		length := d.calculatePieceHashes(index)
		queue <- &PieceJob{
			index,
			hash,
			length,
		}
	}
}

func (d *Downloader) Download() ([]byte, error) {
	log.Println("Starting to download torrent file", d.T.Name)

	queue := make(chan *PieceJob, len(d.T.PieceHashes))
	results := make(chan *PieceResult)

	d.setupJobs(queue, results)

	for _, peer := range d.Peers {
		go d.startDownloadJob(peer, queue, results)
	}

	buf := d.showProgress(results)

	close(queue)
	return buf, nil
}

func (d *Downloader) showProgress(results chan *PieceResult) []byte {
	buf := make([]byte, d.T.Length)
	done := 0
	for done < len(d.T.PieceHashes) {
		res := <-results
		begin, end := d.calculateBoundsForPiece(res.index)
		copy(buf[begin:end], res.buf)
		done++
		percent := float64(done) / float64(len(d.T.PieceHashes)) * 100
		numOfWorkers := runtime.NumGoroutine() - 1

		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numOfWorkers)
	}
	return buf
}

func (d *Downloader) calculatePieceHashes(index int) int {
	begin, end := d.calculateBoundsForPiece(index)
	return end - begin
}

func (d *Downloader) calculateBoundsForPiece(index int) (int, int) {
	begin := index * d.T.PieceLength
	end := begin + d.T.PieceLength
	if end > d.T.Length {
		end = d.T.Length
	}
	return begin, end
}

func (d *Downloader) startDownloadJob(peer peer.Peer, queue chan *PieceJob, results chan *PieceResult) {
	client, err := client.New(peer, d.PeerId, d.T.InfoHash)
	if err != nil {
		log.Printf("Fail to handshake with peer %s.", peer.IP)
		return
	}

	defer client.Conn.Close()

	log.Printf("Completed handshake with %s\n", peer.IP)

	client.SendUnchoke()
	client.SendInterested()

	for job := range queue {
		if !client.Bitfield.HasPiece(job.index) {
			queue <- job
			continue
		}

		buf, err := attemptToDownloadPiece(client, job)
		if err != nil {
			log.Println("Exiting", err)
			queue <- job
			return
		}

		err = checkIntegrity(job, buf)
		if err != nil {
			log.Printf("Piece #%d failed integrity check\n", job.index)
			queue <- job
			continue
		}
		client.SendHave(job.index)
		results <- &PieceResult{index: job.index, buf: buf}
	}

}

func attemptToDownloadPiece(client *client.Client, job *PieceJob) ([]byte, error) {
	state := PieceProgress{
		index:  job.index,
		client: client,
		buf:    make([]byte, job.length),
	}

	client.Conn.SetDeadline(time.Now().Add(utils.PIECE_READ_DEADLINE))

	defer client.Conn.SetDeadline(time.Now())

	for state.downloaded < job.length {
		if !state.client.Choked {
			for state.backlog < utils.MAX_BACKLOG && state.requested < job.length {
				blockSize := utils.MAX_BACKLOG_SIZE
				if job.length-state.requested < blockSize {
					blockSize = job.length - state.requested
				}

				err := client.SendRequest(job.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}

		err := state.readMessage()
		if err != nil {
			return nil, err
		}
	}

	return state.buf, nil
}

func (p *PieceProgress) readMessage() error {
	msg, err := p.client.Read()

	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	switch msg.MessageID {
	case message.MESSAGE_UNCHOKE:
		p.client.Choked = false
	case message.MESSAGE_CHOKE:
		p.client.Choked = true

	case message.MESSAGE_HAVE:
		index, err := message.ParseMessageHave(msg)
		if err != nil {
			return err
		}
		p.client.Bitfield.SetPiece(index)
	case message.MESSAGE_PIECE:
		n, err := message.ParseMessagePiece(p.index, p.buf, msg)
		if err != nil {
			return err
		}
		p.downloaded += n
		p.backlog--
	}

	return nil

}

func checkIntegrity(job *PieceJob, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], job.hash[:]) {
		return fmt.Errorf("Index %d failed integrity check", job.index)
	}
	return nil
}
