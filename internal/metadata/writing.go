package metadata

import (
	"strconv"

	"github.com/bogem/id3v2"
	merrors "github.com/tekofx/musicfixer/internal/errors"
)

func (m Metadata) WriteToFile(filepath string) *merrors.MError {

	// TODO: Add frame parse as option
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return merrors.NewWithArgs(merrors.UnexpectedError, err)
	}

	// Write tags
	tag.AddTextFrame(MetadataTrack, tag.DefaultEncoding(), strconv.Itoa(m.Track))
	if m.Disc != nil {
		tag.AddTextFrame(MetadataDisc, tag.DefaultEncoding(), strconv.Itoa(*m.Disc))
	}
	tag.AddTextFrame(MetadataAlbumArtist, tag.DefaultEncoding(), m.AlbumArtist)

	tag.AddAttachedPicture(m.Picture)

	// Save tag to file
	err = tag.Save()
	if err != nil {
		return merrors.NewWithArgs(merrors.UnexpectedError, err)
	}

	return nil
}
