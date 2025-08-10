package tests

import (
	"log"
	"os"
	"testing"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/metadata"
)

func TestReadMetadata(t *testing.T) {

	t.Run("Read Metadata", testReadMetadata)
	t.Run("Missing metadata", testMissingMetadata)

}

func TestWriteMetadata(t *testing.T) {
	t.Run("Write Metadata", testWriteMetadata)
}

func testReadMetadata(t *testing.T) {
	_, merr, merrors := metadata.GetMetadata("songs/correct_metadata.mp3")
	Assert(t, merr == nil, "Error")
	Assert(t, merrors == nil, "Perror")
}

func testMissingMetadata(t *testing.T) {
	// Missing title
	_, merr, metaErrors := metadata.GetMetadata("songs/missing_title.mp3")
	Assert(t, merr == nil, "Merror not nil")
	Assert(t, metaErrors != nil, "Metaerror is nil")
	Assert(t, metaErrors.Errors[0].Code == merrors.MissingTitle, "Missing title")

	// Missing artist
	_, merr, metaErrors = metadata.GetMetadata("songs/missing_artist.mp3")
	Assert(t, merr == nil, "Merror not nil")
	Assert(t, metaErrors != nil, "Metaerror is nil")
	Assert(t, metaErrors.Errors[0].Code == merrors.MissingArtist, "Missing artist")

	// Missing album
	_, merr, metaErrors = metadata.GetMetadata("songs/missing_album.mp3")
	Assert(t, merr == nil, "Merror not nil")
	Assert(t, metaErrors.Errors[0].Code == merrors.MissingAlbum, "Missing Album")

	// Missing album artist
	_, merr, metaErrors = metadata.GetMetadata("songs/missing_album_artist.mp3")
	Assert(t, merr == nil, "Merror not nil")
	Assert(t, metaErrors.Errors[0].Code == merrors.MissingAlbumArtist, "Missing Album Artist")

	// Missing year
	_, merr, metaErrors = metadata.GetMetadata("songs/missing_year.mp3")
	Assert(t, merr == nil, "Merror not nil")
	Assert(t, metaErrors.Errors[0].Code == merrors.MissingYear, "Missing year")

	// Missing cover
	_, merr, metaErrors = metadata.GetMetadata("songs/missing_cover.mp3")
	Assert(t, merr == nil, "Merror not nil")
	Assert(t, metaErrors.Errors[0].Code == merrors.MissingCover, "Missing Cover")
}

func testWriteMetadata(t *testing.T) {

	tag, err := id3v2.Open("songs/empty_tags.mp3", id3v2.Options{Parse: true})
	if err != nil {
		t.Fatal(err)
	}

	tag.DeleteAllFrames()
	// Write tag to file.mp3
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}

	// Read the image file
	imageData, _ := os.ReadFile("test.jpg")

	m := metadata.Metadata{
		Track: 4,
		Picture: id3v2.PictureFrame{
			Picture:     imageData,
			PictureType: id3v2.PTFrontCover,
			Encoding:    id3v2.EncodingISO,
		},
	}

	merror := m.WriteToFile("songs/empty_tags.mp3")
	Assert(t, merror == nil, "Merror is not nil")

}
