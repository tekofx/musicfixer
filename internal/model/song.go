package model

import "github.com/dhowden/tag"

// Song represents a music file with metadata
type Song struct {
	FilePath    string
	NewFilePath string
	Track       int
	Disc        int
	Title       string
	Picture     tag.Picture
	AlbumName   string
}
