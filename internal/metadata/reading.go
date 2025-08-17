package metadata

import (
	"strconv"
	"strings"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
)

func CheckMetadata(m *id3v2.Tag, path string) []merrors.MError {
	var merrs []merrors.MError

	if m.Artist() == "" {
		merrs = append(merrs, *merrors.New(merrors.MissingArtist, "Missing Artist"))
	}

	if m.Album() == "" {
		merrs = append(merrs, *merrors.New(merrors.MissingAlbum, "Missing Album"))
	}

	if m.Title() == "" {
		merrs = append(merrs, *merrors.New(merrors.MissingTitle, "Missing Title"))
	}

	if GetTrack(m) == -1 {
		merrs = append(merrs, *merrors.New(merrors.MissingTrackNumber, "Missing Track Number"))
	}

	if GetAlbumArtist(m) == "" {
		merrs = append(merrs, *merrors.New(merrors.MissingAlbumArtist, "Missing Album Artist"))
	}

	if m.Year() == "" {
		merrs = append(merrs, *merrors.New(merrors.MissingYear, "Missing Year"))
	}

	if GetPicture(m) == nil {
		merrs = append(merrs, *merrors.New(merrors.MissingCover, "Missing Cover"))
	}

	return merrs

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
	if len(metadata.GetFrames("APIC")) == 0 {
		return nil
	}

	picture := metadata.GetFrames("APIC")[0]
	p := picture.(id3v2.PictureFrame)

	return &p

}

func GetYear(metadata *id3v2.Tag) string {
	return metadata.Year()
}

func GetAlbumArtist(metadata *id3v2.Tag) string {
	return metadata.GetTextFrame("TPE2").Text
}
