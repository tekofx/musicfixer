package model

import (
	"github.com/bogem/id3v2"
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
}
