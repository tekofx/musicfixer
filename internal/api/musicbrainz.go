package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	merrors "github.com/tekofx/musicfixer/internal/errors"
)

type Release struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	ArtistCredit []Artist `json:"artist-credit"`
	Date         *string  `json:"date,omitempty"`
	Country      string   `json:"country,omitempty"`
}

type ArtistDetails struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	SortName string `json:"sort-name"`
}

type Artist struct {
	Name    string        `json:"name"`
	Details ArtistDetails `json:"artist"`
}

type MusicBrainzAlbumResponse struct {
	Created  string     `json:"created"`
	Count    int        `json:"count"`
	Offset   int        `json:"offset"`
	Releases []*Release `json:"releases"`
}

func (r *Release) missingMetadata() bool {
	if r.Date == nil {
		return true
	}

	return false
}

func (m *MusicBrainzAlbumResponse) GetFirstValidRelease() *Release {
	for _, r := range m.Releases {
		if !r.missingMetadata() {
			return r
		}
	}
	return nil
}

func searchAlbum(url string) (*MusicBrainzAlbumResponse, *merrors.MError) {
	res, merr := getRequest(url)
	if merr != nil {
		return nil, merr
	}
	defer res.Body.Close()

	var data MusicBrainzAlbumResponse
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotDecodeJson, err)
	}

	if data.Count == 0 {
		return nil, merrors.New(merrors.NotFound, "Release not found")
	}
	return &data, nil
}

func searchAlbumByNameAndArtist(album string, artist string) (*MusicBrainzAlbumResponse, *merrors.MError) {
	// URL-encode the query properly
	url := fmt.Sprintf("https://musicbrainz.org/ws/2/release?query=artist:%s%%20AND%%20release:%s&fmt=json",
		url.QueryEscape(artist),
		url.QueryEscape(album),
	)

	return searchAlbum(url)
}

func searchAlbumByName(album string) (*MusicBrainzAlbumResponse, *merrors.MError) {
	// URL-encode the query properly
	url := fmt.Sprintf("https://musicbrainz.org/ws/2/release?query=release:%s&fmt=json",
		url.QueryEscape(album),
	)

	return searchAlbum(url)
}

func GetAlbumByNameAndArtist(name string, artist string) (*Release, *merrors.MError) {
	res, merr := searchAlbumByNameAndArtist(name, artist)
	if merr != nil {
		return nil, merr
	}

	return res.GetFirstValidRelease(), nil
}
