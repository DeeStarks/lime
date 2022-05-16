package configs

import "github.com/gdamore/tcell/v2"

type LimeVersion struct {
	Number                 int
	Logo                   []string // Each element represents a line
	Author                 string
	InfoText               string
	DefaultBackgroundColor tcell.Color
	DefaultForegroundColor tcell.Color
	BoxBackgroundColor     tcell.Color
	BoxForegroundColor     tcell.Color
}

// List of versions
var (
	V1 = LimeVersion{
		Number: 1,
		Logo: []string{
			"█░░ █ █▀▄▀█ █▀▀",
			"█▄▄ █ █░▀░█ ██▄",
		},
		Author:                 "DeeStarks",
		InfoText:               "Ctrl+Q = quit | Ctrl+W = write | Ctrl+S = save",
		DefaultBackgroundColor: tcell.Color16,
		DefaultForegroundColor: tcell.ColorReset,
		BoxBackgroundColor:     tcell.ColorReset,
		BoxForegroundColor:     tcell.ColorWhite,
	}
)
