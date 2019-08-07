package manager

//go:generate stringer -type ManageMode

type ManageMode uint8

const (
	Copy ManageMode = iota
	Hardlink
)
