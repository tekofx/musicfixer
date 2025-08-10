package model

import (
	"github.com/tekofx/musicfixer/internal/metadata"
)

// Song represents a music file with metadata
type Song struct {
	FilePath    string
	NewFilePath string
	Track       int
	Disc        *int
	Title       string
	Picture     metadata.Picture
	AlbumName   string
}
