package main

import (
	"flag"
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

func pathExists(path string) (bool, *error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // Path exists
	}
	if os.IsNotExist(err) {
		return false, nil // Path does not exist
	}
	return false, &err // Some other error (e.g., permission denied)
}

func setupFlags() (string, bool) {

	var outputDir string
	outputDir = "output"
	flag.StringVar(&outputDir, "output", "output", "Output directory")

	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "Show changes")
	flag.BoolVar(&dryRun, "d", false, "Show changes")

	// Help flag
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS] [directory]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "FLAGS:\n")
		fmt.Fprintf(os.Stderr, "-d, --dry\t Show changes without renaming\n")
		fmt.Fprintf(os.Stderr, "-o, --output\t Output directory of renamed files\n")
		fmt.Fprintf(os.Stderr, "-h, --help\t Show Help\n")
	}

	flag.Parse()
	return outputDir, dryRun
}

func getDir() string {
	args := flag.Args()

	fmt.Println(args)

	rootDir, _ := os.Getwd()

	if len(args) == 0 {
		return rootDir
	}

	pathExists, err := pathExists(args[0])
	if err != nil {
		log.Fatal(err)
	}
	if !pathExists {
		log.Fatal("Path not exists")
	} else {
		rootDir = args[0]
	}

	return rootDir
}

func main() {
	outputDir, dryRun := setupFlags()
	dir := getDir()

	albumSongs, err := readAlbums(dir)
	if err != nil {
		fmt.Println(err)
	}

	if len(*albumSongs) == 0 {
		fmt.Printf("No mp3 files found in %s", dir)
		os.Exit(0)
	}

	err = renameSongs(*albumSongs, dryRun, outputDir)
	if err != nil {
		fmt.Println(err)
	}

}

func renameSongs(albumSongs map[string]Album, dry bool, outputDir string) *error {

	if !dry {
		os.Mkdir(outputDir, 0700)
	}
	for _, album := range albumSongs {

		outputPath := filepath.Join(outputDir, album.Name)
		if !dry {
			os.Mkdir(outputPath, 0700)
			coverPath := filepath.Join(outputPath, "cover.jpg")
			err := saveCoverArt(album.Songs[0].Metadata, coverPath)
			if err != nil {
				return err
			}
		}

		for _, song := range album.Songs {
			newFilePath := getNewFilePath(song, album)

			if dry {
				fmt.Printf("%s-->%s\n", song.FilePath, newFilePath)
			} else {
				err := os.Rename(song.FilePath, newFilePath)
				if err != nil {
					fmt.Println(err)
					return &err
				}
			}
		}
	}
	return nil
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
