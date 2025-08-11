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
	UnexpectedError    MErrorCode = 0
	WrongFlagsPosition MErrorCode = 1
	FlagNotExists      MErrorCode = 2

	// Metadata 10-99
	MissingTitle               MErrorCode = 10
	MissingArtist              MErrorCode = 11
	MissingTrackNumber         MErrorCode = 12
	MissingAlbum               MErrorCode = 13
	MissingAlbumArtist         MErrorCode = 14
	MissingYear                MErrorCode = 15
	MissingCover               MErrorCode = 16
	CouldNotSaveMetadataToFile MErrorCode = 17

	// FilesystemError 100-199
	MP3FilesNotFound    MErrorCode = 100
	CouldNotOpenFile    MErrorCode = 101
	CouldNotCreateDir   MErrorCode = 102
	CouldNotCreateFile  MErrorCode = 103
	CouldNotWriteToFile MErrorCode = 104
	CouldNotRenameFile  MErrorCode = 105
	CouldNotDeleteDirs  MErrorCode = 106
	PathNotExists       MErrorCode = 107

	// Api Requests 200-299
	CouldNotFetchImage         MErrorCode = 200
	UnexpectecStatusCode       MErrorCode = 201
	CouldNotCopyResponseToFile MErrorCode = 202
	CouldNotCreateRequest      MErrorCode = 203
	CouldNotGetResponse        MErrorCode = 204
	CouldNotDecodeJson         MErrorCode = 205
)
