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

	albumSongs, merror, errors := model.ReadAlbums(dir)

	if merror != nil {
		fmt.Printf("merror: %v\n", merror)
	}

	if errors != nil {
		fmt.Println("Error reading songs metadata:")
		for _, error := range errors {
			fmt.Printf("%v\n", error)
		}
		os.Exit(0)
	}

	if albumSongs == nil {
		fmt.Printf("No songs found in %s", dir)
		os.Exit(0)
	}

	model.SetNewFilePaths(albumSongs)

	if dry {
		flags.DryRun(albumSongs, outputDir)
	}

	err := model.RenameSongs(*albumSongs, outputDir)
	if err != nil {
		fmt.Printf("Error renaming songs: %v\n", err)
		os.Exit(0)
	}

	fmt.Printf("Done! All song renamed in %s", dir)

	if removeOriginalFolder {
		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Printf("Error removing original directories: %v", err)
		}
	}

}
