package tests

import (
	"testing"

	"github.com/tekofx/musicfixer/internal/api"
)

func TestMusicBrainz(t *testing.T) {
	t.Run("Get album", searchAlbum)
	t.Run("Download cover", downloadCover)

}

func searchAlbum(t *testing.T) {
	album, merr := api.SearchAlbum("siames", "bounce")
	AssertMErrorNotNil(t, merr)
	Assert(t, album.Releases[0].Title == "BOUNCE INTO THE MUSIC", "Album not corresponds")
}

func downloadCover(t *testing.T) {
	merr := api.SaveReleaseCover("18638c82-5408-483b-b3f2-027a1c4cdd4f")
	AssertMErrorNotNil(t, merr)
	RemoveFile("test.jpg")
}
