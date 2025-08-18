package model

import (
	"fmt"
	"strings"

	"github.com/tekofx/musicfixer/internal/api"
	merrors "github.com/tekofx/musicfixer/internal/errors"
)

type Album struct {
	Name      string
	Artist    string
	Year      string
	Songs     []Song
	MultiDisk bool
}

func (a *Album) AddSong(song Song) {
	if len(a.Songs) == 0 {
		a.Name = song.AlbumName
		a.Artist = song.AlbumArtist
		a.Year = song.Year
	}
	a.Songs = append(a.Songs, song)
}

func (a *Album) FixMetadata() *merrors.MError {
	meta, merr := api.GetAlbumByNameAndArtist(a.Name, a.Songs[0].Artist)
	if merr != nil {
		return merr
	}
	year := strings.Split(*meta.Date, "-")[0]
	cover, merr := api.GetReleaseCover(meta.ID)
	if merr != nil {
		return merr
	}

	artistName := meta.ArtistCredit[0].Details.SortName
	if strings.Contains(artistName, ",") {
		artistSplit := strings.Split(meta.ArtistCredit[0].Details.SortName, ", ")
		artistName = fmt.Sprintf("%s %s", artistSplit[1], artistSplit[0])
	}

	fmt.Println(year)

	for _, s := range a.Songs {
		if s.AlbumArtist == "" {
			s.AlbumArtist = artistName
		}

		if s.Year != year {
			s.Year = year
		}

		if s.Picture == nil {
			s.AddCover(cover)
		}

		s.UpdateFile()

	}

	return nil

}
