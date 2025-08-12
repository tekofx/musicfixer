package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
	Year        string
	Picture     *id3v2.PictureFrame
	AlbumName   string
	Artist      string
	MErrors     []merrors.MError
}

func NewSong(filepath string) (*Song, *merrors.MError) {
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotOpenFile, "Error while opening mp3 file:", err)
	}
	defer tag.Close()

	songMetadataErrors := metadata.CheckMetadata(tag, filepath)

	return &Song{
		FilePath:  filepath,
		Title:     tag.Title(),
		AlbumName: tag.Album(),
		Year:      tag.Year(),
		Artist:    metadata.GetAlbumArtist(tag),
		Track:     metadata.GetTrack(tag),
		Disc:      metadata.GetDisc(tag),
		Picture:   metadata.GetPicture(tag),
		MErrors:   songMetadataErrors,
	}, nil

}

func (song *Song) UpdateFile() *merrors.MError {
	// TODO: Add frame parse as option
	tag, err := id3v2.Open(song.FilePath, id3v2.Options{Parse: true})
	if err != nil {
		return merrors.NewWithArgs(merrors.UnexpectedError, err)
	}

	// Write tags
	tag.AddTextFrame(metadata.MetadataTrack, tag.DefaultEncoding(), strconv.Itoa(song.Track))
	if song.Disc != nil {
		tag.AddTextFrame(metadata.MetadataDisc, tag.DefaultEncoding(), strconv.Itoa(*song.Disc))
	}
	tag.AddTextFrame(metadata.MetadataAlbumArtist, tag.DefaultEncoding(), song.Artist)

	tag.AddAttachedPicture(*song.Picture)

	// Save tag to file
	err = tag.Save()
	if err != nil {
		return merrors.NewWithArgs(merrors.UnexpectedError, err)
	}
	return nil
}

func (song *Song) ContainsMetadataError(code merrors.MErrorCode) bool {
	for _, m2 := range song.MErrors {
		if m2.Code == code {
			return true
		}
	}
	return false
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
