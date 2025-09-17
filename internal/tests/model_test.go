package tests

import (
	"os"
	"testing"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func TestCreateModels(t *testing.T) {
	t.Run("New Song", testNewSong)
	t.Run("New Song Fails", testNewSongFails)

	os.Mkdir("test", 0755)
	t.Run("Read Album", testReadAlbum)
	t.Run("Read Album Fails", testReadAlbumFails)

	os.RemoveAll("test")

}
func testNewSong(t *testing.T) {
	_, merr := model.NewSong("files/correct_metadata.mp3")
	AssertMErrorNotNil(t, merr)

}
func testNewSongFails(t *testing.T) {
	_, merr := model.NewSong("a.mp3")
	AssertMError(t, merr, merrors.CouldNotOpenFile, "Error while opening mp3 file: open a.mp3: no such file or directory")
}

func testReadAlbumFails(t *testing.T) {
	musicCollection := model.NewMusicCollection()
	merr := musicCollection.ReadAlbums("test")
	AssertMError(t, merr, merrors.MP3FilesNotFound, "Not mp3 files found in test")
}

func TestAddSong(t *testing.T) {
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

	Assert(t, album.Name == "Bounce Into the Music", "Album name not corresponds")
	Assert(t, album.Artist == "SIAMES", "Album artist not corresponds")
	Assert(t, album.Year == "2016", "Year not corresponds")
	Assert(t, len(album.Songs) == 1, "Song not added")
	Assert(t, !album.MultiDisk, "Album is multidisk")
}
