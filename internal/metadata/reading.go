package metadata

import (
	"strconv"
	"strings"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
)

type Metadata struct {
	Title       string
	Album       string
	AlbumArtist string
	Year        string
	Track       int
	Disc        *int
	Picture     id3v2.PictureFrame
}

const (
	MetadataTrack       = "TRCK"
	MetadataDisc        = "TPOS"
	MetadataCover       = "APIC"
	MetadataAlbumArtist = "TPE2"
)

func checkMetadata(m *id3v2.Tag, path string) *merrors.SongMetadataError {
	songMetadataErrors := merrors.SongMetadataError{
		SongPath: path,
	}

	if m.Artist() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingArtist, "Missing Artist"))
	}

	if m.Album() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingAlbum, "Missing Album"))
	}

	if m.Title() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingTitle, "Missing Title"))
	}

	if getTrack(m) == -1 {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingTrackNumber, "Missing Track Number"))
	}

	if getAlbumArtist(m) == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingAlbumArtist, "Missing Album Artist"))
	}

	if m.Year() == "" {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingYear, "Missing Year"))

	}

	if len(m.GetFrames("APIC")) == 0 {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingCover, "Missing Cover"))
	}

	if len(songMetadataErrors.Errors) > 0 {
		return &songMetadataErrors
	}
	return nil

}

func GetMetadata(path string) (*Metadata, *merrors.MError, *merrors.SongMetadataError) {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return nil, merrors.NewWithArgs(merrors.CouldNotOpenFile, "Error while opening mp3 file:", err), nil
	}
	defer tag.Close()

	songMetadataErrors := checkMetadata(tag, path)
	if songMetadataErrors != nil {
		return nil, nil, songMetadataErrors
	}

	metadata := Metadata{
		Title:       tag.Title(),
		Album:       tag.Album(),
		Year:        tag.Year(),
		AlbumArtist: getAlbumArtist(tag),
		Track:       getTrack(tag),
		Disc:        getDisc(tag),
		Picture:     getPicture(tag),
	}

	return &metadata, nil, nil

}

func getTrack(metadata *id3v2.Tag) int {

	str := strings.Split(metadata.GetTextFrame("TRCK").Text, "/")[0]

	if str == "" {
		return -1
	}

	num, _ := strconv.Atoi(str)

	return num
}

func getDisc(metadata *id3v2.Tag) *int {
	str := strings.Split(metadata.GetTextFrame("TPOS").Text, "/")[0]

	if str == "" {
		return nil
	}

	num, _ := strconv.Atoi(str)

	return &num

}

func getPicture(metadata *id3v2.Tag) id3v2.PictureFrame {
	picture := metadata.GetFrames("APIC")[0]
	p := picture.(id3v2.PictureFrame)

	return p

}

func getAlbumArtist(metadata *id3v2.Tag) string {
	return metadata.GetTextFrame("TPE2").Text
}
