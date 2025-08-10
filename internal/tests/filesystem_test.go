package tests

import (
	"os"
	"testing"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func tearDown() {
	os.RemoveAll("test")
}

func setup() {
	os.Mkdir("test", 0755)
}

func TestFilesystem(t *testing.T) {
	setup()

	// Filesystem
	t.Run("Open file", testReadAlbum)

	tearDown()

}

func testReadAlbum(t *testing.T) {
	albums, merror, metaErrors := model.ReadAlbums("./")
	Assert(t, albums == nil, "Albumn not nil")
	Assert(t, len(metaErrors) == 6, "Metaerrors len != 6")

	albums, merror, metaErrors = model.ReadAlbums("test")
	Assert(t, albums == nil, "Albumn not nil")
	Assert(t, len(metaErrors) == 0, "Metaerrors len != 6")
	AssertMError(t, merror, merrors.MP3FilesNotFound, "Not mp3 files found in test")
}
