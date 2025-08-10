package merrors

import (
	"fmt"
	"strings"
)

type SongMetadataError struct {
	SongPath string
	Errors   []MError
}

func (e SongMetadataError) Error() string {
	msg := ""
	for _, error := range e.Errors {
		msg += fmt.Sprintf(" - %s\n", error.Message)
	}

	return fmt.Sprintf("Song %s misses: \n%s", e.SongPath, msg)
}

type MError struct {
	Code    MErrorCode
	Message string
}

func (m MError) Print() {
	fmt.Println(m.Message)
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

func NewWithArgs(code MErrorCode, messages ...any) *MError {
	var parts []string
	for _, arg := range messages {
		parts = append(parts, fmt.Sprint(arg))
	}
	return &MError{
		Code:    code,
		Message: strings.Join(parts, " "),
	}
}

type MErrorCode int

const (

	// Common errors
	UnexpectedError = 0

	// Missing Metadata 10-99
	MissingTitle       = 10
	MissingArtist      = 11
	MissingTrackNumber = 12
	MissingAlbum       = 13
	MissingAlbumArtist = 14
	MissingYear        = 15
	MissingCover       = 16

	// FilesystemError 100-199
	MP3FilesNotFound    = 100
	CouldNotOpenFile    = 101
	CouldNotCreateDir   = 102
	CouldNotCreateFile  = 103
	CouldNotWriteToFile = 104
	CouldNotRenameFile  = 105
	CouldNotDeleteDirs  = 106
	PathNotExists       = 107
)
