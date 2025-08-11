package api

import (
	"encoding/json"
	"fmt"
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

	res, merr := GetRequest(url)
	if merr != nil {
		return nil, merr
	}
	defer res.Body.Close()

	var data model.MusicBrainzAlbumResponse
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotDecodeJson, err)
	}
	return &data, nil
}
