package perrors

import (
	"fmt"
)

type SongMetadataError struct {
	SongPath string
	Errors   []string
}

func (e SongMetadataError) Error() string {
	msg := ""
	for _, error := range e.Errors {
		msg += fmt.Sprintf(" - %s\n", error)
	}

	return fmt.Sprintf("Song %s misses: \n%s", e.SongPath, msg)
}
