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

	addCover(tag, m.Picture)

	// Save tag to file
	err = tag.Save()
	if err != nil {
		return merrors.NewWithArgs(merrors.UnexpectedError, err)
	}

	return nil
}

func addCover(tag *id3v2.Tag, picture Picture) *merrors.MError {

	pictureFrame := id3v2.PictureFrame{
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Encoding:    id3v2.EncodingUTF8,
		Picture:     picture.Data,
	}

	// Add the PictureFrame to the tag
	tag.AddAttachedPicture(pictureFrame)
	return nil
}
