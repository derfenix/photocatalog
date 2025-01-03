package metadata

import (
	"fmt"
	"io"

	"github.com/rwcarlsen/goexif/exif"
)

type Exif struct{}

func (j Exif) Extract(_ string, data io.Reader) (Metadata, error) {
	decode, err := exif.Decode(data)
	if err != nil {
		return Metadata{}, fmt.Errorf("decode exif: %w", err)
	}

	meta := Metadata{}

	created, err := decode.DateTime()
	if err != nil {
		return Metadata{}, fmt.Errorf("parse datetime: %w", err)
	}

	meta.Created = created

	return meta, nil
}
