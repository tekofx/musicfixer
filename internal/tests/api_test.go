package tests

import (
	"testing"

	"github.com/tekofx/musicfixer/internal/api"
)

func TestMusicBrainz(t *testing.T) {
	t.Run("Get album", getAlbum)

}

func getAlbum(t *testing.T) {
	api.SearchAlbum("siames", "bounce")
}
