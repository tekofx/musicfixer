package main

import (
	"fmt"
	"os"

	"github.com/tekofx/musicfixer/internal/flags"
	"github.com/tekofx/musicfixer/internal/model"
)

func main() {
	outputDir, dry, removeOriginalFolder := flags.SetupFlags()
	dir := flags.GetDir()

	albumSongs, err := model.ReadAlbums(dir)
	if err != nil {
		fmt.Printf("Error reading songs: %v\n", err)
		os.Exit(0)
	}

	model.SetNewFilePaths(albumSongs)

	if dry {
		flags.DryRun(albumSongs, outputDir)
	}

	err = model.RenameSongs(*albumSongs, outputDir)
	if err != nil {
		fmt.Printf("Error renaming songs: %v\n", err)
		os.Exit(0)
	}

	if removeOriginalFolder {
		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Printf("Error removing original directories: %v", err)
		}
	}

}
