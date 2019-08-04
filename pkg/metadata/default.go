package metadata

import (
	"path"
	"strings"
	"time"
)

const defaultTimeLayout = "20060102_150405"

// DefaultExtractor extract metadata from all file types, not covered by special extractors
//
// Gets the meta data from the file's name
type DefaultExtractor struct {
	Layout string
}

// NewDefaultExtractor returns new DefaultExtractor's instance
func NewDefaultExtractor() *DefaultExtractor {
	return &DefaultExtractor{Layout: defaultTimeLayout}
}

// NewDefaultExtractorWithLayout returns DefaultExtractor with custom time layout
func NewDefaultExtractorWithLayout(l string) *DefaultExtractor {
	return &DefaultExtractor{Layout: l}
}

// Extract returns Metadata from specified filename using its name to parse Time
func (d *DefaultExtractor) Extract(fp string) (Metadata, error) {
	_, fName := path.Split(fp)

	// Remove extension
	fName = strings.Replace(fName, path.Ext(fName), "", 1)

	// If there more than one photo in one second, cameras append ~N to the end of file name (before extension)
	if strings.ContainsRune(fName, '~') {
		fName = fName[:strings.IndexRune(fName, '~')]
	}

	t, err := time.ParseInLocation(d.Layout, fName, time.Local)
	if err != nil {
		return Metadata{}, err
	}
	return Metadata{Time: t}, nil
}
