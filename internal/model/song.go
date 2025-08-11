package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/metadata"
	"github.com/tekofx/musicfixer/internal/utils"
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

func NewSong(filepath string) (*Song, *merrors.MError, *merrors.SongMetadataError) {
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

func (song *Song) SaveCover(outputFilePath string) *merrors.MError {
	var err error

	// Check if cover art exists
	if song.Picture != nil {
		return merrors.New(merrors.MissingCover, "No cover found to save")

	}

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotCreateFile, "Failed to create cover file:", err)
	}
	defer outputFile.Close()

	// Write the cover art data to the file
	_, err = outputFile.Write(song.Picture.Picture)
	if err != nil {
		err = fmt.Errorf("failed to write cover art: %w", err)
		return merrors.NewWithArgs(merrors.CouldNotWriteToFile, "Failed to write cover to file:", err)
	}

	return nil
}

func (song *Song) SetNewFilePath(album Album) string {
	track := song.Track
	var newName string
	var trackString string

	if track < 10 {
		trackString = fmt.Sprintf("0%d", track)
	} else {
		trackString = fmt.Sprintf("%d", track)
	}

	if album.MultiDisk {
		disc := song.Disc
		newName = fmt.Sprintf("Disc %d - %s. %s.mp3", disc, trackString, utils.CleanFilename(song.Title))

	} else {
		newName = fmt.Sprintf("%s. %s.mp3", trackString, utils.CleanFilename(song.Title))
	}

	return filepath.Join("output", utils.CleanFilename(album.Name), newName)
}
