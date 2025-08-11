package tests

import (
	"testing"

	"github.com/tekofx/musicfixer/internal/api"
)

func TestMusicBrainz(t *testing.T) {
	t.Run("Get album", getAlbum)

}

func getAlbum(t *testing.T) {
	album, merr := api.SearchAlbum("siames", "bounce")

	AssertMErrorNotNil(t, merr)
	Assert(t, album.Releases[0].Title == "BOUNCE INTO THE MUSIC", "Album not corresponds")
}
