package metadata

import (
	"fmt"
	"log"
	"os"

	"github.com/dhowden/tag"
)

func checkMetadata(m tag.Metadata, path string) error {

	if m.Album() == "" {
		return fmt.Errorf("Song %s does not have album", path)
	}

	if m.Title() == "" {
		return fmt.Errorf("Song %s does not have title", path)
	}

	track, _ := m.Track()
	disc, _ := m.Disc()

	if track == 0 {
		return fmt.Errorf("Song %s does not have track number", path)
	}
	if disc == 0 {
		return fmt.Errorf("Song %s does not have disc number", path)
	}

	return nil

}

func GetMetadata(path string) (tag.Metadata, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", path, err)
		return nil, err
	}
	defer file.Close()
	m, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = checkMetadata(m, path)
	if err != nil {
		return nil, err
	}

	return m, nil

}
