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
	_, merr, _ := metadata.GetMetadata("songs/correct_metadata.mp3")
	Assert(t, merr == nil, "Could not open file")

	_, merr, _ = metadata.GetMetadata("a.mp3")
	AssertMError(t, merr, merrors.CouldNotOpenFile, "Error while opening mp3 file: open a.mp3: no such file or directory")

}
func testReadAlbum(t *testing.T) {
	albums, merror, metaErrors := model.ReadAlbums("./")
	Assert(t, albums == nil, "Albumn not nil")
	Assert(t, len(metaErrors) == 7, "Metaerrors len != 7")

	albums, merror, metaErrors = model.ReadAlbums("test")
	Assert(t, albums == nil, "Albumn not nil")
	Assert(t, len(metaErrors) == 0, "Metaerrors len != 6")
	AssertMError(t, merror, merrors.MP3FilesNotFound, "Not mp3 files found in test")
}
