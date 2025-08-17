package metadata

import "github.com/bogem/id3v2"

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
	MetadataYear        = "TDRC"
)
