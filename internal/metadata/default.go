package metadata

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

const timeFormat = "20060102_150405"

var (
	ErrBadPrefix     = errors.New("bad prefix")
	ErrBadNameFormat = errors.New("bad name format")
)

type Default struct {
	TimeFormat string
	Prefix     string
}

func (d *Default) Extract(fp string, _ io.Reader) (Metadata, error) {
	format := d.TimeFormat
	if format == "" {
		format = timeFormat
	}

	if d.Prefix != "" && !strings.HasPrefix(fp, d.Prefix) {
		return Metadata{}, fmt.Errorf("%w: expect a prefix %s, got %s", ErrBadPrefix, d.Prefix, fp)
	}

	fp = filepath.Base(fp)

	leftLimit := len(d.Prefix)
	rightLimit := leftLimit + len(format)

	if len(fp) < rightLimit {
		return Metadata{}, fmt.Errorf("%w: too short", ErrBadNameFormat)
	}

	created, err := time.Parse(format, fp[leftLimit:rightLimit])
	if err != nil {
		return Metadata{}, fmt.Errorf("parse time: %w", err)
	}

	meta := Metadata{
		Created: created,
	}

	return meta, nil
}
