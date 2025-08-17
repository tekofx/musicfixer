package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
)

// Check a condition. If the condition is false, the assert fails and shows failMessage
func Assert(t *testing.T, predicate bool, failMessage string) {
	if !predicate {
		fmt.Println("Test failed:", failMessage)
		t.FailNow()
	}
}

func removeMetadataFromFile(filepath string) *merrors.MError {
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return merrors.NewWithArgs(merrors.CouldNotOpenFile, err)
	}

	tag.DeleteAllFrames()
	if err = tag.Save(); err != nil {
		return merrors.NewWithArgs(merrors.CouldNotSaveMetadataToFile, err)
	}
	return nil
}

func AssertMErrorNotNil(t *testing.T, error *merrors.MError) {
	if error != nil {
		fmt.Printf("Test failed\n%s\n", error.Message)
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

func RemoveFile(filepath string) {
	os.Remove(filepath)
}
