package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

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

func setupFlags() (string, bool, bool) {

	var outputDir string
	outputDir = "output"
	flag.StringVar(&outputDir, "output", "output", "Output directory")

	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "Show changes")
	flag.BoolVar(&dryRun, "d", false, "Show changes")

	var removeOriginalFolder bool
	flag.BoolVar(&removeOriginalFolder, "remove", false, "Remove original folder")
	flag.BoolVar(&removeOriginalFolder, "r", false, "Remove original folder")
	// Help flag
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS] [directory]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "FLAGS:\n")
		fmt.Fprintf(os.Stderr, "-d, --dry\t Show changes without renaming\n")
		fmt.Fprintf(os.Stderr, "-o, --output\t Output directory of renamed files\n")
		fmt.Fprintf(os.Stderr, "-h, --help\t Show Help\n")
		fmt.Fprintf(os.Stderr, "-r, --remove\t Remove original folder\n")

	}

	flag.Parse()
	return outputDir, dryRun, removeOriginalFolder
}

func getDir() string {
	args := flag.Args()
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
