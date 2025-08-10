package metadata

import (
	"log"
	"strings"

	"github.com/bogem/id3v2"
	perrors "github.com/tekofx/musicfixer/internal/errors"
)

type Metadata struct {
	Title   string
	Album   string
	Year    string
	Track   string
	Disc    string
	Picture Picture
}

type Picture struct {
	Data      []byte
	Extension string
}

func checkMetadata(m *id3v2.Tag, path string) *perrors.SongMetadataError {
	songMetadataErrors := perrors.SongMetadataError{
		SongPath: path,
	}

	if m.Album() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Album Name")
	}

	if m.Title() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Song title")
	}

	if getTrack(m) == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Track number")
	}

	if len(m.GetFrames("APIC")) == 0 {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, "Cover")
	}

	if len(songMetadataErrors.Errors) > 0 {
		return &songMetadataErrors
	}
	return nil

}

func GetMetadata(path string) (*Metadata, error, *perrors.SongMetadataError) {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	songMetadataErrors := checkMetadata(tag, path)
	if songMetadataErrors != nil {
		return nil, nil, songMetadataErrors
	}

	metadata := Metadata{
		Title:   tag.Title(),
		Album:   tag.Album(),
		Year:    tag.Year(),
		Track:   getTrack(tag),
		Disc:    getDisc(tag),
		Picture: getPicture(tag),
	}

	return &metadata, nil, nil

}

func getTrack(metadata *id3v2.Tag) string {
	return strings.Split(metadata.GetTextFrame("TRCK").Text, "/")[0]
}

func getDisc(metadata *id3v2.Tag) string {
	return strings.Split(metadata.GetTextFrame("TPOS").Text, "/")[0]
}

func getPicture(metadata *id3v2.Tag) Picture {
	picture := metadata.GetFrames("APIC")[0]
	p := picture.(id3v2.PictureFrame)

	ext := ".jpg"
	if p.MimeType == "image/png" {
		ext = ".png"
	}

	return Picture{
		Data:      p.Picture,
		Extension: ext,
	}

}
