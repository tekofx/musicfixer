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

func TestWriteMetadata(t *testing.T) {
	t.Run("Write Metadata", testWriteMetadata)
	merr := removeMetadataFromFile("files/empty_tags.mp3")
	AssertMErrorNotNil(t, merr)
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

	// Update song file
	merr = song.UpdateFile()
	AssertMErrorNotNil(t, merr)

	// Check metadata has been written correctly
	song, merr = model.NewSong("files/empty_tags.mp3")

	AssertMErrorNotNil(t, merr)
	Assert(t, len(song.MErrors) == 6, "Not all metadata errors")

	Assert(t, song.ContainsMetadataError(merrors.MissingTitle), "Not missing title")
	Assert(t, song.ContainsMetadataError(merrors.MissingArtist), "Not missing artist")
	Assert(t, song.ContainsMetadataError(merrors.MissingAlbum), "Not missing album")
	Assert(t, song.ContainsMetadataError(merrors.MissingTrackNumber), "Not missing track number")
	Assert(t, song.ContainsMetadataError(merrors.MissingAlbumArtist), "Not missing album artist")
	Assert(t, song.ContainsMetadataError(merrors.MissingYear), "Not missing year")
}

func TestFixMetadata(t *testing.T) {
	t.Run("Fixing Metadata Wrong Song", testFixMetadataWrongSong)
	t.Run("Fixing Metadata Correct Song", testFixMetadataCorrectSong)

}

func testFixMetadataCorrectSong(t *testing.T) {
	album := model.Album{
		Name:   "",
		Artist: "",
		Year:   "",
	}
	album.AddSong(model.Song{
		Title:       "The Wolf",
		Artist:      "SIAMES",
		AlbumArtist: "SIAMES",
		AlbumName:   "Bounce Into the Music",
		Year:        "2016",
	})

	merr := album.FixMetadata()
	AssertMErrorNotNil(t, merr)
}

func testFixMetadataWrongSong(t *testing.T) {
	album := model.Album{
		Name:   "",
		Artist: "",
		Year:   "",
	}
	album.AddSong(model.Song{})
	merr := album.FixMetadata()
	AssertMError(t, merr, merrors.NotFound, "404")
}
