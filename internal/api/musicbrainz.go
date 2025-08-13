package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	merrors "github.com/tekofx/musicfixer/internal/errors"
)

func SearchAlbum(artist string, album string) (*MusicBrainzAlbumResponse, *merrors.MError) {
	// URL-encode the query properly
	url := fmt.Sprintf("https://musicbrainz.org/ws/2/release?query=artist:%s%%20AND%%20release:%s&fmt=json",
		url.QueryEscape(artist),
		url.QueryEscape(album),
	)

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
		return nil, merrors.New(merrors.EmptyResponse, "Album response is empty")
	}
	return &data, nil
}
