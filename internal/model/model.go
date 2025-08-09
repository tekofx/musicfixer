package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrz1836/go-sanitize"
	perrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/metadata"
)

func getSong(path string) (*Song, error, *perrors.SongMetadataError) {

	m, err, songMetadataErrors := metadata.GetMetadata(path)
	if err != nil {
		return nil, err, nil
	}

	if songMetadataErrors != nil {

		return nil, nil, songMetadataErrors
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

	return &song, nil, nil
}

func SetNewFilePaths(albumSongs *map[string]Album) {
	for _, album := range *albumSongs {
		for i := range album.Songs {
			newFilePath := setNewFilePath(album.Songs[i], album)
			album.Songs[i].NewFilePath = newFilePath
		}
	}
}

func setNewFilePath(song Song, album Album) string {
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
		newName = fmt.Sprintf("Disc %d - %s. %s.mp3", disc, trackString, sanitize.AlphaNumeric(song.Title, false))
	} else {
		newName = fmt.Sprintf("%s. %s.mp3", trackString, song.Title)
	}

	return filepath.Join("output", album.Name, newName)
}

func RenameSongs(albumSongs map[string]Album, outputDir string) error {
	os.Mkdir(outputDir, 0700)
	for _, album := range albumSongs {
		outputPath := filepath.Join(outputDir, sanitize.AlphaNumeric(album.Name, false))
		os.Mkdir(outputPath, 0700)
		coverPath := filepath.Join(outputPath, "cover.jpg")
		err := saveCover(album.Songs[0], coverPath)
		if err != nil {
			return err
		}

		for _, song := range album.Songs {
			err := os.Rename(song.FilePath, song.NewFilePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ReadAlbums(searchDir string) (*map[string]Album, []perrors.SongMetadataError) {
	// Initialize a map to group songs by album
	albumSongs := make(map[string]Album)

	var perrors []perrors.SongMetadataError

	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file is an MP3
		if filepath.Ext(path) != ".mp3" {
			return nil
		}

		song, err, songMetadataErrors := getSong(path)
		if err != nil {
			return nil
		}

		if songMetadataErrors != nil {
			perrors = append(perrors, *songMetadataErrors)
			return nil
		}

		album := albumSongs[song.AlbumName]
		album.Name = song.AlbumName
		album.Songs = append(album.Songs, *song)
		if !album.MultiDisk {
			if song.Disc > 1 {
				album.MultiDisk = true
			}
		}

		// Add the song to the appropriate album group
		albumSongs[song.AlbumName] = album

		return nil
	})

	if err != nil {
		return nil, nil
	}

	if len(perrors) > 0 {
		return nil, perrors
	}

	if len(albumSongs) == 0 {
		err = fmt.Errorf("Not mp3 files found in %s", searchDir)
		return nil, nil
	}

	return &albumSongs, nil
}

func saveCover(song Song, outputFilePath string) error {
	var err error

	// Retrieve the cover art data
	songPicture := song.Picture

	// Check if cover art exists
	if len(songPicture.Data) == 0 {
		err = fmt.Errorf("no cover art found")

		return err
	}

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		err = fmt.Errorf("failed to create output file: %w", err)
		return err
	}
	defer outputFile.Close()

	// Write the cover art data to the file
	_, err = outputFile.Write(songPicture.Data)
	if err != nil {
		err = fmt.Errorf("failed to write cover art: %w", err)
		return err
	}

	return nil
}
