package flags

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func pathExists(path string) (bool, *merrors.MError) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // Path exists
	}
	if os.IsNotExist(err) {
		return false, nil // Path does not exist
	}
	return false, merrors.NewWithArgs(merrors.UnexpectedError, err) // Some other error (e.g., permission denied)
}

func SetupFlags() (string, bool, bool, bool, *merrors.MError) {
	var outputDir string
	outputDir = "output"
	flag.StringVar(&outputDir, "output", "output", "Output directory")
	flag.StringVar(&outputDir, "o", "output", "Output directory")

	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "Show changes")
	flag.BoolVar(&dryRun, "d", false, "Show changes")

	var removeOriginalFolder bool
	flag.BoolVar(&removeOriginalFolder, "remove", false, "Remove original folder")
	flag.BoolVar(&removeOriginalFolder, "r", false, "Remove original folder")

	var completeMetadata bool
	flag.BoolVar(&completeMetadata, "fix", false, "Complete missing metadata")
	flag.BoolVar(&completeMetadata, "f", false, "Complete missing metadata")

	// Help flag
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS] [directory]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "FLAGS:\n")
		fmt.Fprintf(os.Stderr, "-d, --dry\t Show changes without renaming\n")
		fmt.Fprintf(os.Stderr, "-o, --output\t Output directory of renamed files\n")
		fmt.Fprintf(os.Stderr, "-h, --help\t Show Help\n")
		fmt.Fprintf(os.Stderr, "-r, --remove\t Remove original folder\n")
		fmt.Fprintf(os.Stderr, "-f, --fix\t Completes missing metadata\n")
	}

	flag.Parse()

	if len(flag.Args()) > 1 {
		return "", false, false, false, merrors.NewWithArgs(merrors.WrongFlagsPosition, "Flags go before directory: musicfixer", strings.Join(flag.Args()[1:], " "), flag.Arg(0))
	}

	return outputDir, dryRun, removeOriginalFolder, completeMetadata, nil
}

func GetDir() (*string, *merrors.MError) {
	args := flag.Args()
	rootDir, _ := os.Getwd()

	if len(args) == 0 {
		return &rootDir, nil
	}

	pathExists, merr := pathExists(args[0])
	if merr != nil {
		return nil, merr
	}
	if !pathExists {
		return nil, merrors.New(merrors.PathNotExists, "Path not exists")
	} else {
		rootDir = args[0]
	}

	return &rootDir, nil
}

func DryRun(musicCollection *model.MusicCollection, outputDir string) {
	if musicCollection.HasMetaErrors() {
		musicCollection.PrintMetaErrors()
		os.Exit(0)
	}
	for _, album := range musicCollection.Albums {
		outputPath := filepath.Join(outputDir, album.Name)
		coverPath := filepath.Join(outputPath, "cover.jpg")
		fmt.Printf("Cover: %s\n", coverPath)
		for _, song := range album.Songs {
			fmt.Printf("%s-->%s\n", song.FilePath, song.NewFilePath)
		}
	}
	os.Exit(0)
}
