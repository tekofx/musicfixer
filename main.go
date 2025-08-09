package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
)

func main() {
	outputDir, dry, removeOriginalFolder := setupFlags()
	dir := getDir()

	albumSongs, err := readAlbums(dir)
	if err != nil {
		fmt.Println(*err)
		os.Exit(0)
	}

	getNewFilePaths(albumSongs)

	if dry {
		dryRun(albumSongs, outputDir)
	}

	err = renameFiles(*albumSongs, outputDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if removeOriginalFolder {
		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func dryRun(albumSongs *map[string]Album, outputDir string) {
	for _, album := range *albumSongs {
		outputPath := filepath.Join(outputDir, album.Name)
		coverPath := filepath.Join(outputPath, "cover.jpg")
		fmt.Printf("Cover: %s\n", coverPath)
		for _, song := range album.Songs {
			fmt.Printf("%s-->%s\n", song.FilePath, song.NewFilePath)
		}
	}
	os.Exit(0)
}

func renameFiles(albumSongs map[string]Album, outputDir string) *error {
	os.Mkdir(outputDir, 0700)
	for _, album := range albumSongs {
		outputPath := filepath.Join(outputDir, album.Name)
		os.Mkdir(outputPath, 0700)
		coverPath := filepath.Join(outputPath, "cover.jpg")
		err := saveCoverArt(album.Songs[0].Metadata, coverPath)
		if err != nil {
			return err
		}

		for _, song := range album.Songs {
			err := os.Rename(song.FilePath, song.NewFilePath)
			if err != nil {
				return &err
			}
		}
	}
	return nil
}

func getNewFilePaths(albumSongs *map[string]Album) {
	for _, album := range *albumSongs {
		for i := range album.Songs {
			newFilePath := getNewFilePath(album.Songs[i], album)
			album.Songs[i].NewFilePath = newFilePath
		}
	}
}

func getNewFilePath(song Song, album Album) string {
	track, _ := song.Metadata.Track()
	var newName string
	var trackString string

	if track < 10 {
		trackString = fmt.Sprintf("0%d", track)
	} else {
		trackString = fmt.Sprintf("%d", track)
	}

	if album.MultiDisk {
		disc, _ := song.Metadata.Disc()
		newName = fmt.Sprintf("Disc %d - %s. %s.mp3", disc, trackString, song.Metadata.Title())
	} else {
		newName = fmt.Sprintf("%s. %s.mp3", trackString, song.Metadata.Title())
	}

	return filepath.Join("output", album.Name, newName)
}

func readAlbums(rootDir string) (*map[string]Album, *error) {
	// Initialize a map to group songs by album
	albumSongs := make(map[string]Album)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file is an MP3
		if filepath.Ext(path) != ".mp3" {
			err = fmt.Errorf("Not mp3 files found in %s", path)
			return err
		}

		file, err := os.Open(path)
		if err != nil {

			log.Printf("Error opening file %s: %v\n", path, err)
			return err
		}
		defer file.Close()

		m, err := tag.ReadFrom(file)
		if err != nil {
			log.Fatal(err)
		}

		album := m.Album()

		// Create a new song entry
		song := Song{
			FilePath: path,
			Metadata: m,
		}

		albumArray := albumSongs[album]
		albumArray.Name = album
		albumArray.Songs = append(albumArray.Songs, song)
		if !albumArray.MultiDisk {
			disc, _ := song.Metadata.Disc()
			if disc > 1 {
				albumArray.MultiDisk = true
			}
		}

		// Add the song to the appropriate album group
		albumSongs[album] = albumArray

		return nil
	})
	if err != nil {
		return nil, &err
	}
	return &albumSongs, nil
}

func saveCoverArt(m tag.Metadata, outputFilePath string) *error {
	var err error

	// Retrieve the cover art data
	coverData := m.Picture()

	// Check if cover art exists
	if len(coverData.Data) == 0 {
		err = fmt.Errorf("no cover art found")

		return &err
	}

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		err = fmt.Errorf("failed to create output file: %w", err)
		return &err
	}
	defer outputFile.Close()

	// Write the cover art data to the file
	_, err = outputFile.Write(coverData.Data)
	if err != nil {
		err = fmt.Errorf("failed to write cover art: %w", err)
		return &err
	}

	return nil
}
