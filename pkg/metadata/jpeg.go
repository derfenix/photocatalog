package metadata

import (
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

// JpegExtractor meta data extractor for the jpeg files
type JpegExtractor struct {
}

// NewJpegExtractor returns new JpegExtractor
func NewJpegExtractor() *JpegExtractor {
	return &JpegExtractor{}
}

// Extract returns Metadata from specified jpeg file reading its exif data
//
// TODO: Fallback to default extractor on exif reading/parsing error
func (j *JpegExtractor) Extract(fp string) (Metadata, error) {
	f, err := os.Open(fp)
	if err != nil {
		return Metadata{}, err
	}
	x, err := exif.Decode(f)
	if err != nil {
		return Metadata{}, err
	}

	time, err := x.DateTime()
	if err != nil {
		return Metadata{}, err
	}
	meta := Metadata{
		Time: time,
	}

	return meta, nil
}
