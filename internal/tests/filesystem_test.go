package tests

import (
	"os"
	"testing"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/metadata"
	"github.com/tekofx/musicfixer/internal/model"
)

func TestFilesystem(t *testing.T) {
	t.Run("Open file", testOpenFile)

	os.Mkdir("test", 0755)
	t.Run("Read Album", testReadAlbum)
	os.RemoveAll("test")

}
func testOpenFile(t *testing.T) {
	_, merr, _ := metadata.ReadMetadata("files/correct_metadata.mp3")
	AssertMErrorNotNil(t, merr)

	_, merr, _ = metadata.ReadMetadata("a.mp3")
	AssertMError(t, merr, merrors.CouldNotOpenFile, "Error while opening mp3 file: open a.mp3: no such file or directory")

}
func testReadAlbum(t *testing.T) {
	musicCollection := model.NewMusicCollection()

	merr, metaErrors := musicCollection.ReadAlbums("./")
	AssertMErrorNotNil(t, merr)
	Assert(t, len(metaErrors) == 7, "Metaerrors len != 7")

	musicCollection = model.NewMusicCollection()

	merr, metaErrors = musicCollection.ReadAlbums("test")
	Assert(t, len(metaErrors) == 0, "Metaerrors len != 6")
	AssertMError(t, merr, merrors.MP3FilesNotFound, "Not mp3 files found in test")
}
