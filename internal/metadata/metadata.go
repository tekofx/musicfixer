package metadata

import (
	"log"
	"os"

	"github.com/dhowden/tag"
	perrors "github.com/tekofx/musicfixer/internal/errors"
)

func checkMetadata(m tag.Metadata, path string) *perrors.SongMetadataError {
	songMetadataErrors := perrors.SongMetadataError{
		SongPath: path,
	}

	if m.Album() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Album Name")

	}

	if m.Title() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Song title")
	}

	if m.Picture() == nil {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Cover")
	}

	track, _ := m.Track()

	if track == 0 {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Track number")
	}

	if len(songMetadataErrors.Errors) > 0 {
		return &songMetadataErrors
	}
	return nil

}

func GetMetadata(path string) (tag.Metadata, error, *perrors.SongMetadataError) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", path, err)
		return nil, err, nil
	}
	defer file.Close()
	m, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
		return nil, err, nil
	}

	songMetadataErrors := checkMetadata(m, path)
	if songMetadataErrors != nil {
		return nil, nil, songMetadataErrors
	}

	return m, nil, nil

}
