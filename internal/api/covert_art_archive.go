package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func SaveReleaseCover(releaseId string) *merrors.MError {
	url := fmt.Sprintf("http://coverartarchive.org/release/%s",
		url.QueryEscape(releaseId),
	)

	res, merr := GetRequest(url)
	if merr != nil {
		return merr
	}
	defer res.Body.Close()

	var data model.CoverArtResponse
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotDecodeJson, err)
	}

	data.Images[0].Save("test.jpg")
	return nil

}
