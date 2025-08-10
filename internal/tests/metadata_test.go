package tests

import (
	"fmt"
	"testing"

	merrors "github.com/tekofx/musicfixer/internal/errors"
	"github.com/tekofx/musicfixer/internal/metadata"
)

func TestReadMetadata(t *testing.T) {

	// Filesystem
	t.Run("Open file", testOpenFile)

	t.Run("Get Correct Metadata", testGetCorrectMetadata)
	t.Run("Missing metadata", testGetIncorrectMetadata)

}

func testOpenFile(t *testing.T) {
	_, merr, _ := metadata.GetMetadata("correct_metadata.mp3")
	Assert(t, merr == nil, "Could not open file")

	_, merr, _ = metadata.GetMetadata("a.mp3")
	Assert(t, merr.Code == merrors.CouldNotOpenFile, "Could not open file")
}

func testGetCorrectMetadata(t *testing.T) {
	_, err, merrors := metadata.GetMetadata("correct_metadata.mp3")
	fmt.Println(err)
	Assert(t, err == nil, "Error")
	Assert(t, merrors == nil, "Perror")
}

func testGetIncorrectMetadata(t *testing.T) {
	_, err, merrors := metadata.GetMetadata("missing_title.mp3")

	Assert(t, err == nil, "Error")
	Assert(t, merrors != nil, "Perror")
}
