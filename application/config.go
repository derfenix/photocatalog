package application

import (
	"fmt"
	"slices"
)

type Mode string

const (
	ModeCopy     Mode = "copy"
	ModeHardlink Mode = "hardlink"
	ModeSymlink  Mode = "symlink"
	ModeMove     Mode = "move"
)

var SupportedModes = []Mode{ModeHardlink, ModeSymlink, ModeMove, ModeCopy}

type Config struct {
	SourceDir    string
	TargetDir    string
	Mode         Mode
	Overwrite    bool
	DirMode      uint64
	FileMode     uint64
	Watch        bool
	SkipFullSync bool
}

func (c *Config) Validate() error {
	if c.SourceDir == "" {
		return fmt.Errorf("source dir is required")
	}

	if c.TargetDir == "" {
		return fmt.Errorf("target dir is required")
	}

	if !slices.Contains(SupportedModes, c.Mode) {
		return fmt.Errorf("invalid mode %s, supported modes: %s", c.Mode, SupportedModes)
	}

	if c.SkipFullSync && !c.Watch {
		return fmt.Errorf("skip full sync and watch disabled — nothing to do")
	}

	return nil
}
