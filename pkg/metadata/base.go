package metadata

import (
	"time"
)

// Metadata contains meta data for the files have to be processed
type Metadata struct {
	Time time.Time
}

// Extractor interface for Metadata extractors
type Extractor interface {
	Extract(string) (Metadata, error)
}
