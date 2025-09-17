package tests

import (
	"os"
	"testing"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func TestFilesystem(t *testing.T) {
	t.Run("Open file", testOpenFile)

	os.Mkdir("test", 0755)
	t.Run("Read Album", testReadAlbum)
	os.RemoveAll("test")

}
func testOpenFile(t *testing.T) {
	_, merr := model.NewSong("files/correct_metadata.mp3")
	AssertMErrorNotNil(t, merr)

	_, merr = model.NewSong("a.mp3")
	AssertMError(t, merr, merrors.CouldNotOpenFile, "Error while opening mp3 file: open a.mp3: no such file or directory")

}
func testReadAlbum(t *testing.T) {
	musicCollection := model.NewMusicCollection()
	merr := musicCollection.ReadAlbums("./")
	AssertMErrorNotNil(t, merr)

	musicCollection = model.NewMusicCollection()
	merr = musicCollection.ReadAlbums("test")
	AssertMError(t, merr, merrors.MP3FilesNotFound, "Not mp3 files found in test")
}
