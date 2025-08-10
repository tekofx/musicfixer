package tests

import (
	"testing"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/metadata"
)

func TestReadMetadata(t *testing.T) {

	// Filesystem
	t.Run("Open file", testOpenFile)

	t.Run("Get Correct Metadata", testGetCorrectMetadata)
	t.Run("Missing metadata", testMissingMetadata)

}

func testOpenFile(t *testing.T) {
	_, merr, _ := metadata.GetMetadata("songs/correct_metadata.mp3")
	Assert(t, merr == nil, "Could not open file")

	_, merr, _ = metadata.GetMetadata("a.mp3")
	AssertMError(t, merr, merrors.CouldNotOpenFile, "Error while opening mp3 file: open a.mp3: no such file or directory")

}

func testGetCorrectMetadata(t *testing.T) {
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
