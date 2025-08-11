package song

import (
	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/metadata"
)

// Song represents a music file with metadata
type Song struct {
	FilePath    string
	NewFilePath string
	Track       int
	Disc        *int
	Title       string
	Picture     *id3v2.PictureFrame
	AlbumName   string
	Artist      string
}

func New(filepath string) (*Song, *merrors.MError, *merrors.SongMetadataError) {
	metadata, merror, metaerrors := metadata.ReadMetadata(filepath)

	if merror != nil {
		return nil, merror, nil
	}

	if metaerrors != nil {
		return nil, nil, metaerrors
	}

	return &Song{
		FilePath:  filepath,
		Title:     metadata.Title,
		AlbumName: metadata.Album,
		Artist:    metadata.AlbumArtist,
		Track:     metadata.Track,
		Disc:      metadata.Disc,
	}, nil, nil

}
