package metadata

import (
	"strconv"
	"strings"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
)

func CheckMetadata(m *id3v2.Tag, path string) *merrors.SongMetadataError {
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

	if GetTrack(m) == -1 {
		songMetadataErrors.Errors = append(songMetadataErrors.Errors, *merrors.New(merrors.MissingTrackNumber, "Missing Track Number"))
	}

	if GetAlbumArtist(m) == "" {
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

func GetTrack(metadata *id3v2.Tag) int {

	str := strings.Split(metadata.GetTextFrame("TRCK").Text, "/")[0]

	if str == "" {
		return -1
	}

	num, _ := strconv.Atoi(str)

	return num
}

func GetDisc(metadata *id3v2.Tag) *int {
	str := strings.Split(metadata.GetTextFrame("TPOS").Text, "/")[0]

	if str == "" {
		return nil
	}

	num, _ := strconv.Atoi(str)

	return &num

}

func GetPicture(metadata *id3v2.Tag) *id3v2.PictureFrame {
	picture := metadata.GetFrames("APIC")[0]
	p := picture.(id3v2.PictureFrame)

	return &p

}

func GetAlbumArtist(metadata *id3v2.Tag) string {
	return metadata.GetTextFrame("TPE2").Text
}
