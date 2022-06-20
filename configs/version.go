package configs

import "github.com/gdamore/tcell/v2"

type LimeVersion struct {
	Number                 int
	DefaultBackgroundColor tcell.Color
	DefaultForegroundColor tcell.Color
	BoxBackgroundColor     tcell.Color
	BoxForegroundColor     tcell.Color
}

// List of versions
var (
	V1 = LimeVersion{
		Number: 1,
		DefaultBackgroundColor: tcell.Color16,
		DefaultForegroundColor: tcell.ColorReset,
		BoxBackgroundColor:     tcell.ColorReset,
		BoxForegroundColor:     tcell.ColorWhite,
	}
)
