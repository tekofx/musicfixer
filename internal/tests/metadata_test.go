package tests

import (
	"os"
	"testing"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func TestReadMetadata(t *testing.T) {

	t.Run("Read Metadata", testReadMetadata)
	t.Run("Missing metadata", testMissingMetadata)

}

func TestWriteMetadata(t *testing.T) {
	t.Run("Write Metadata", testWriteMetadata)
}

func testReadMetadata(t *testing.T) {
	song, merr := model.NewSong("files/correct_metadata.mp3")
	AssertMErrorNotNil(t, merr)
	Assert(t, len(song.MErrors) == 0, "Missing title")

}

func testMissingMetadata(t *testing.T) {
	// Missing title
	song, merr := model.NewSong("files/missing_title.mp3")
	AssertMErrorNotNil(t, merr)
	Assert(t, len(song.MErrors) == 1, "No metadata errors")
	Assert(t, song.MErrors[0].Code == merrors.MissingTitle, "Missing title")

	// Missing artist
	song, merr = model.NewSong("files/missing_artist.mp3")
	AssertMErrorNotNil(t, merr)
	Assert(t, len(song.MErrors) == 1, "No metadata errors")
	Assert(t, song.MErrors[0].Code == merrors.MissingArtist, "Missing artist")

	// Missing album
	song, merr = model.NewSong("files/missing_album.mp3")
	AssertMErrorNotNil(t, merr)
	Assert(t, len(song.MErrors) == 1, "No metadata errors")
	Assert(t, song.MErrors[0].Code == merrors.MissingAlbum, "Missing album")

	// Missing album artist
	song, merr = model.NewSong("files/missing_album_artist.mp3")
	Assert(t, len(song.MErrors) == 1, "No metadata errors")
	Assert(t, song.MErrors[0].Code == merrors.MissingAlbumArtist, "Missing album artist")

	// Missing year
	song, merr = model.NewSong("files/missing_year.mp3")
	Assert(t, len(song.MErrors) == 1, "No metadata errors")
	Assert(t, song.MErrors[0].Code == merrors.MissingYear, "Missing year")

	// Missing cover
	song, merr = model.NewSong("files/missing_cover.mp3")
	Assert(t, len(song.MErrors) == 1, "No metadata errors")
	Assert(t, song.MErrors[0].Code == merrors.MissingCover, "Missing cover")

}

func testWriteMetadata(t *testing.T) {
	// Read the image file
	imageData, _ := os.ReadFile("files/cover.jpg")

	song, merr := model.NewSong("files/empty_tags.mp3")
	AssertMErrorNotNil(t, merr)

	song.Picture = &id3v2.PictureFrame{
		Picture:     imageData,
		PictureType: id3v2.PTFrontCover,
		Encoding:    id3v2.EncodingISO,
	}

	merr = song.UpdateFile()
	AssertMErrorNotNil(t, merr)

	// Check metadata has been written correctly
	song, merr = model.NewSong("files/empty_tags.mp3")
	AssertMErrorNotNil(t, merr)
	Assert(t, len(song.MErrors) == 5, "No metadata errors")

	Assert(t, song.ContainsMetadataError(merrors.MissingCover), "Missing cover")
	Assert(t, song.ContainsMetadataError(merrors.MissingTitle), "Missing title")
	Assert(t, song.ContainsMetadataError(merrors.MissingAlbum), "Missing album")
	Assert(t, song.ContainsMetadataError(merrors.MissingAlbumArtist), "Missing album artist")
	Assert(t, song.ContainsMetadataError(merrors.MissingYear), "Missing year")
	removeMetadataFromFile("files/empty_tags.mp3")
}
