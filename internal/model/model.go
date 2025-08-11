package model

import (
	"fmt"
	"os"
	"path/filepath"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/song"
	"github.com/tekofx/musicfixer/internal/utils"
)

func SetNewFilePaths(musicCollection *MusicCollection) {
	for _, album := range musicCollection.Albums {
		for i := range album.Songs {
			newFilePath := setNewFilePath(album.Songs[i], album)
			album.Songs[i].NewFilePath = newFilePath
		}
	}
}

func setNewFilePath(song song.Song, album Album) string {
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

func RenameSongs(musicCollection MusicCollection, outputDir string) *merrors.MError {
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
		err := saveCover(album.Songs[0], coverPath)
		if err != nil {
			return err
		}

		for _, song := range album.Songs {

			err := os.Rename(song.FilePath, song.NewFilePath)
			if err != nil {
				return merrors.NewWithArgs(merrors.CouldNotRenameFile, "Could not rename file", song.FilePath, song.NewFilePath, err)
			}
		}
	}
	return nil
}

func ReadAlbums(searchDir string) (*MusicCollection, *merrors.MError, []merrors.SongMetadataError) {
	// Initialize a map to group songs by album
	musicCollection := MusicCollection{}

	var perrors []merrors.SongMetadataError
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

		song, err2, songMetadataErrors := song.New(path)
		if err2 != nil {
			merr = err2
			return filepath.SkipAll
		}

		if songMetadataErrors != nil {
			perrors = append(perrors, *songMetadataErrors)
			return nil
		}

		// Add song to musiccollection
		musicCollection.AddSong(*song)

		return nil
	})

	if merr != nil {
		return nil, merr, nil
	}

	if len(perrors) > 0 {
		return nil, nil, perrors
	}

	if len(musicCollection.Albums) == 0 {
		return nil, merrors.NewWithArgs(merrors.MP3FilesNotFound, "Not mp3 files found in", searchDir), nil
	}

	return &musicCollection, nil, nil
}

func saveCover(song song.Song, outputFilePath string) *merrors.MError {
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
