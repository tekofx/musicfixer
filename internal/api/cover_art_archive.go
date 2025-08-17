package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	merrors "github.com/tekofx/musicfixer/internal/errors"
)

type CoverArtResponse struct {
	Images []CoverArtImage `json:images`
}

type CoverArtImage struct {
	Id    int    `json:id`
	Front bool   `json:front`
	Image string `json:image`
}

func getCoverArt(releaseId string) (*CoverArtResponse, *merrors.MError) {
	url := fmt.Sprintf("http://coverartarchive.org/release/%s",
		url.QueryEscape(releaseId),
	)

	res, merr := getRequest(url)
	if merr != nil {
		return nil, merr
	}
	defer res.Body.Close()

	var data CoverArtResponse
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotDecodeJson, err)
	}
	return &data, nil
}

func GetReleaseCover(releaseId string) ([]byte, *merrors.MError) {
	data, merr := getCoverArt(releaseId)
	if merr != nil {
		return nil, merr
	}

	imageData, merr := data.Images[0].toBytes()
	if merr != nil {
		return nil, merr
	}

	return imageData, nil

}

func SaveReleaseCover(releaseId string) *merrors.MError {
	data, merr := getCoverArt(releaseId)
	if merr != nil {
		return merr
	}

	data.Images[0].save("test.jpg")
	return nil

}
func (c *CoverArtImage) toBytes() ([]byte, *merrors.MError) {
	// 1. Fetch the image
	resp, err := http.Get(c.Image)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotFetchImage, err)
	}
	defer resp.Body.Close()

	// 2. Check that the response is OK
	if resp.StatusCode != http.StatusOK {
		return nil, merrors.NewWithArgs(merrors.UnexpectecStatusCode, resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotGetImageBytes, err)
	}
	return imageData, nil
}

func (c *CoverArtImage) save(filepath string) *merrors.MError {
	// 1. Fetch the image
	resp, err := http.Get(c.Image)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotFetchImage, err)
	}
	defer resp.Body.Close()

	// 2. Check that the response is OK
	if resp.StatusCode != http.StatusOK {
		return merrors.NewWithArgs(merrors.UnexpectecStatusCode, resp.StatusCode)
	}

	// 3. Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotCreateFile, err)
	}
	defer file.Close()

	// 4. Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {

		return merrors.NewWithArgs(merrors.CouldNotCopyResponseToFile, err)
	}
	return nil
}
