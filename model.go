package main

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

type Album struct {
	Name      string
	Songs     []Song
	MultiDisk bool
}

func getSong(path string) (*Song, error) {

	m, err := getMetadata(path)
	if err != nil {
		return nil, err
	}

	track, _ := m.Track()
	disc, _ := m.Disc()
	song := Song{
		FilePath:  path,
		Title:     m.Title(),
		Track:     track,
		Disc:      disc,
		Picture:   *m.Picture(),
		AlbumName: m.Album(),
	}

	return &song, nil
}
