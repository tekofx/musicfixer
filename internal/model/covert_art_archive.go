package model

import (
	"io"
	"net/http"
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

func (c *CoverArtImage) Save(filepath string) *merrors.MError {
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
