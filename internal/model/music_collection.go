package model

import (
	"fmt"
	"os"
	"path/filepath"

	merrors "github.com/tekofx/musicfixer/internal/errors"
)

type MusicCollection struct {
	Albums map[string]Album
}

func NewMusicCollection() *MusicCollection {
	return &MusicCollection{
		Albums: make(map[string]Album), // No need to take address
	}
}

func (mc *MusicCollection) HasMetaErrors() bool {
	for _, album := range mc.Albums {
		for _, song := range album.Songs {
			if len(song.MErrors) > 0 {
				return true
			}
		}
	}
	return false
}

func (mc *MusicCollection) FixMetadata() *merrors.MError {
	for _, album := range mc.Albums {
		merr := album.FixMetadata()
		if merr != nil {
			return merr
		}

	}

	return nil
}

func (mc *MusicCollection) PrintMetaErrors() {
	for _, album := range mc.Albums {
		for _, song := range album.Songs {
			if len(song.MErrors) > 0 {
				fmt.Printf("%s:\n", song.FilePath)
				for _, error := range song.MErrors {
					fmt.Printf(" - %v\n", error.Message)
				}
			}
		}
	}
}

func (mc *MusicCollection) AddSong(song Song) {
	album := mc.Albums[song.AlbumName]
	album.AddSong(song)
	if !album.MultiDisk {
		if song.Disc != nil && *song.Disc > 1 {
			album.MultiDisk = true
		}
	}
	mc.Albums[song.AlbumName] = album
}

func (mc *MusicCollection) AddAlbum(album Album) {
	mc.Albums[album.Name] = album
}

func (m *MusicCollection) SetNewFilePaths(outputDir string) {
	for _, album := range m.Albums {
		for i := range album.Songs {
			album.Songs[i].SetNewFilePath(album, outputDir)
		}
	}
}

func (musicCollection *MusicCollection) RenameSongs(outputDir string) *merrors.MError {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotCreateDir, "Could not create dir", outputDir, err)
	}
	for _, album := range musicCollection.Albums {
		outputPath := filepath.Dir(album.Songs[0].NewFilePath)
		err = os.MkdirAll(outputPath, 0755)
		if err != nil {
			return merrors.NewWithArgs(merrors.CouldNotCreateDir, "Could not create dir", outputPath, err)
		}
		coverPath := filepath.Join(outputPath, "cover.jpg")
		err := album.Songs[0].SaveCover(coverPath)
		if err != nil {
			return err
		}

		for _, song := range album.Songs {

			fmt.Println(song.Title)

			err := os.Rename(song.FilePath, song.NewFilePath)
			if err != nil {
				return merrors.NewWithArgs(merrors.CouldNotRenameFile, "Could not rename file", song.FilePath, song.NewFilePath, err)
			}
		}
	}
	return nil
}

func (m *MusicCollection) ReadAlbums(searchDir string) *merrors.MError {
	// Initialize a map to group songs by album
	var merr *merrors.MError

	filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file is an MP3
		if filepath.Ext(path) != ".mp3" {
			return nil
		}

		song, err2 := NewSong(path)
		if err2 != nil {
			merr = err2
			return filepath.SkipAll
		}

		// Add song to musiccollection
		m.AddSong(*song)

		return nil
	})

	if merr != nil {
		return merr
	}

	if len(m.Albums) == 0 {
		return merrors.NewWithArgs(merrors.MP3FilesNotFound, "Not mp3 files found in", searchDir)
	}

	return nil
}
