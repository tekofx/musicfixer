package main

import "github.com/dhowden/tag"

// Song represents a music file with metadata
type Song struct {
	FilePath    string
	Metadata    tag.Metadata
	NewFilePath string
}

type Album struct {
	Name      string
	Songs     []Song
	MultiDisk bool
}
