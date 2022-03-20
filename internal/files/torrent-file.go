package files

import "strings"

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
	IsHTTP      bool
}

func NewTorrentFile(path string) (TorrentFile, error) {

	bto, err := ParseMetaInfo(path)
	infohash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	pieceHashes, err := bto.Info.splitHashes()
	if err != nil {
		return TorrentFile{}, nil
	}

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infohash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
		IsHTTP:      strings.Contains(bto.Announce, "http"),
	}

	return t, nil

}
