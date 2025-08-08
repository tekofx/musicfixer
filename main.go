package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
)

// Song represents a music file with metadata
type Song struct {
	FilePath string
	Metadata tag.Metadata
}

type Album struct {
	Name      string
	Songs     []Song
	MultiDisk bool
}

func main() {
	// Define the root directory to search (you can change this)
	rootDir := "./" // Current directory

	// Initialize a map to group songs by album
	albumSongs, err := readAlbums(rootDir)
	if err != nil {
		fmt.Println(err)
	}

	renameSongs(*albumSongs)

	//saveCoverArt(m, "cover.jpg")
}

func renameSongs(albumSongs map[string]Album) {
	for key, album := range albumSongs {
		fmt.Printf("Album: %s\n", key)
		fmt.Printf("Num songs: %d\n", len(album.Songs))

		for _, song := range album.Songs {
			fmt.Printf("  File: %s\n", song.FilePath)
			fmt.Printf("  Title: %q\n", song.Metadata.Title())

			dir := filepath.Dir(song.FilePath)
			track, _ := song.Metadata.Track()

			var newName string
			if album.MultiDisk {
				disc, _ := song.Metadata.Disc()
				newName = fmt.Sprintf("Disc %d - %d. %s", disc, track, song.Metadata.Title())
			} else {
				newName = fmt.Sprintf("%d. %s", track, song.Metadata.Title())
			}

			newFilepath := filepath.Join(dir, newName)
			fmt.Printf("New name: %s\n", newFilepath)

			//os.Rename(song.FilePath, newFilepath)
		}
	}
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
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file %s: %v\n", path, err)
			return nil
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
		fmt.Println(err)
		return nil, &err
	}
	return &albumSongs, nil
}

func saveCoverArt(m tag.Metadata, outputFilePath string) error {
	// Retrieve the cover art data
	coverData := m.Picture()

	// Check if cover art exists
	if len(coverData.Data) == 0 {
		return fmt.Errorf("no cover art found")
	}

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Write the cover art data to the file
	_, err = outputFile.Write(coverData.Data)
	if err != nil {
		return fmt.Errorf("failed to write cover art: %w", err)
	}

	fmt.Printf("Cover art saved to: %s\n", outputFilePath)
	return nil
}
