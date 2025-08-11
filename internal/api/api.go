package api

import (
	"net/http"

	merrors "github.com/tekofx/musicfixer/internal/errors"
)

func GetRequest(url string) (*http.Response, *merrors.MError) {

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

	if res.StatusCode != http.StatusOK {
		return nil, merrors.NewWithArgs(merrors.UnexpectecStatusCode, res.StatusCode)
	}

	return res, nil
}
