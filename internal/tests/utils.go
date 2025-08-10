package tests

import (
	"fmt"
	"os"
	"testing"

	merrors "github.com/tekofx/musicfixer/internal/errors"
)

// Check a condition. If the condition is false, the assert fails and shows failMessage
func Assert(t *testing.T, predicate bool, failMessage string) {
	if !predicate {
		fmt.Println("Test failed:", failMessage)
		t.FailNow()
	}
}

func AssertMError(t *testing.T, error *merrors.MError, code merrors.MErrorCode, message string) {
	if nil == error {
		fmt.Println("Test failed because error is empty.")
		t.FailNow()
	}
	Assert(t, error.Code == code && error.Message == message, fmt.Sprintf("\n[%d - %s] \nwas expected but \n[%d - %s] \nwas found\n", code, message, error.Code, error.Message))
}

func generateMp3File() string {
	// Create a blank (silent) MP3 file
	file, err := os.Create("blank.mp3")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write minimal silent MP3 data (example: 1-second silence, ~1KB)
	// This is a simplified placeholder. Use real MP3 bytes.
	// You can get this from a real silent MP3 or generate via external tool.
	silentMP3 := []byte{
		0x49, 0x44, 0x33, 0x03, 0x00, 0x00, 0x00, 0x00, // ID3 header (optional)
		// Actual MP3 frames would follow here
		// For real use, include valid MPEG audio frames
	}

	_, err = file.Write(silentMP3)
	if err != nil {
		panic(err)
	}

	return "blank.mp3"
}
