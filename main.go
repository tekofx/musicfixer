package main

import (
	"fmt"
	"os"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/flags"
	"github.com/tekofx/musicfixer/internal/model"
)

func main() {
	outputDir, dry, removeOriginalFolder, completeMetadata, merr := flags.SetupFlags()
	if merr != nil {
		merr.Print()
		os.Exit(0)
	}
	dir, merr := flags.GetDir()
	if merr != nil {
		merr.Print()
		os.Exit(0)
	}

	musicCollection := model.NewMusicCollection()

	merr = musicCollection.ReadAlbums(*dir)

	if merr != nil {
		merr.Print()
		os.Exit(0)
	}

	if dry {
		flags.DryRun(musicCollection, outputDir)
		os.Exit(0)
	}

	if completeMetadata {
		merr = musicCollection.FixMetadata()
		if merr != nil {
			merr.Print()
			os.Exit(0)
		}
	}

	if musicCollection.HasMetaErrors() {
		musicCollection.PrintMetaErrors()
		os.Exit(0)
	}

	musicCollection.SetNewFilePaths(outputDir)

	merr = musicCollection.RenameSongs(outputDir)
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
