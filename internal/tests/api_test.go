package tests

import (
	"testing"

	"github.com/tekofx/musicfixer/internal/api"
	merrors "github.com/tekofx/musicfixer/internal/errors"
)

func TestMusicBrainz(t *testing.T) {
	t.Run("Get album", searchAlbumByNameAndArtist)
	t.Run("Download cover", downloadCover)
	t.Run("Get release cover", getReleaseCover)

}

func searchAlbumByNameAndArtist(t *testing.T) {
	album, merr := api.GetAlbumByNameAndArtist("bounce", "siames")
	AssertMErrorNotNil(t, merr)
	Assert(t, album.Title == "BOUNCE INTO THE MUSIC", "Album not corresponds")

	album, merr = api.GetAlbumByNameAndArtist("siames", "bounce")
	AssertMError(t, merr, merrors.NotFound, "Release not found")
}

func downloadCover(t *testing.T) {
	merr := api.SaveReleaseCover("18638c82-5408-483b-b3f2-027a1c4cdd4f")
	AssertMErrorNotNil(t, merr)
	RemoveFile("test.jpg")

	merr = api.SaveReleaseCover("a")
	AssertMError(t, merr, merrors.NotFound, "400")
}

func getReleaseCover(t *testing.T) {
	_, merr := api.GetReleaseCover("18638c82-5408-483b-b3f2-027a1c4cdd4f")
	AssertMErrorNotNil(t, merr)

	_, merr = api.GetReleaseCover("a")
	AssertMError(t, merr, merrors.NotFound, "400")

}
