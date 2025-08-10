package merrors

import (
	"fmt"
)

type SongMetadataError struct {
	SongPath string
	Errors   []MError
}

func (e SongMetadataError) Error() string {
	msg := ""
	for _, error := range e.Errors {
		msg += fmt.Sprintf(" - %s\n", error)
	}

	return fmt.Sprintf("Song %s misses: \n%s", e.SongPath, msg)
}

type MError struct {
	Code    MErrorCode
	Message string
}

func Unexpected(message string) *MError {
	return &MError{
		Code:    UnexpectedError,
		Message: message,
	}
}

func New(code MErrorCode, message string) *MError {
	return &MError{
		Code:    code,
		Message: message,
	}
}

type MErrorCode int

const (

	// Common errors
	UnexpectedError = 0

	// Missing Metadata
	MissingTitle       = 10
	MissingArtist      = 11
	MissingTrackNumber = 12
	MissingAlbum       = 13
	MissingAlbumArtist = 14
	MissingYear        = 15
	MissingCover       = 16
)
