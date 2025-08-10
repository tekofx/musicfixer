package main

import (
	"fmt"
	"os"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/flags"
	"github.com/tekofx/musicfixer/internal/model"
)

func main() {
	outputDir, dry, removeOriginalFolder, merr := flags.SetupFlags()
	if merr != nil {
		merr.Print()
		os.Exit(0)
	}
	dir, merr := flags.GetDir()
	if merr != nil {
		merr.Print()
		os.Exit(0)
	}

	albumSongs, merr, errors := model.ReadAlbums(*dir)

	if merr != nil {
		merr.Print()
		os.Exit(0)
	}

	if errors != nil {
		fmt.Println("Error reading songs metadata:")
		for _, error := range errors {
			fmt.Printf("%v\n", error)
		}
		os.Exit(0)
	}

	model.SetNewFilePaths(albumSongs)

	if dry {
		flags.DryRun(albumSongs, outputDir)
	}

	merr = model.RenameSongs(*albumSongs, outputDir)
	if merr != nil {
		merr.Print()
		os.Exit(0)
	}

	fmt.Printf("Done! All song renamed in %s", outputDir)

	if removeOriginalFolder {
		err := os.RemoveAll(*dir)
		if err != nil {
			merrors.NewWithArgs(merrors.CouldNotDeleteDirs, "Error removing original directories:", err).Print()
		}
	}

}
