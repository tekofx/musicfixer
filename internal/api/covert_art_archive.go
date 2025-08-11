package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/model"
)

func SaveReleaseCover(releaseId string) *merrors.MError {
	url := fmt.Sprintf("http://coverartarchive.org/release/%s",
		url.QueryEscape(releaseId),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotCreateRequest, err)
	}

	// Request JSON format
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MyMusicApp/1.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotGetResponse, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return merrors.NewWithArgs(merrors.UnexpectecStatusCode, res.StatusCode)
	}

	var data model.CoverArtResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotDecodeJson, err)
	}

	fmt.Println(data)

	data.Images[0].Save("test.jpg")
	return nil

}
