package tests

import (
	"fmt"
	"testing"

	"github.com/tekofx/musicfixer/internal/metadata"
)

func TestReadMetadata(t *testing.T) {
	t.Run("Get Correct Metadata", testGetCorrectMetadata)
	t.Run("Missing metadata", testGetIncorrectMetadata)

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
