package model

import (
	"fmt"

	"github.com/tekofx/musicfixer/internal/song"
)

type Album struct {
	Name      string
	Songs     []song.Song
	MultiDisk bool
}

func (a Album) AddSong(song song.Song) {
	a.Songs = append(a.Songs, song)
}

type MusicCollection struct {
	Albums *map[string]Album
}

func NewMusicCollection() *MusicCollection {
	albums := make(map[string]Album)
	return &MusicCollection{
		Albums: &albums,
	}
}

func (mc MusicCollection) AddSong(song song.Song) {
	album := mc.Albums[song.AlbumName]
	fmt.Println(album)
	album.AddSong(song)
	if !album.MultiDisk {
		if song.Disc != nil && *song.Disc > 1 {
			album.MultiDisk = true
		}
	}
	mc.Albums[album.Name] = album
}

func (mc MusicCollection) AddAlbum(album Album) {
	mc.Albums[album.Name] = album
}
