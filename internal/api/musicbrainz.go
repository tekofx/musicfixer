package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func SearchAlbum(artist string, album string) (*model.MusicBrainzAlbumResponse, *merrors.MError) {
	// URL-encode the query properly
	url := fmt.Sprintf("http://musicbrainz.org/ws/2/release?query=artist:%s%%20AND%%20release:%s&fmt=json",
		url.QueryEscape(artist),
		url.QueryEscape(album),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotCreateRequest, err)
	}

	// Request JSON format
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MyMusicApp/1.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotGetResponse, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, merrors.NewWithArgs(merrors.UnexpectecStatusCode, res.StatusCode)
	}

	var data model.MusicBrainzAlbumResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotDecodeJson, err)
	}
	return &data, nil
}
