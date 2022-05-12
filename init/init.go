package init

import (
	"github.com/DeeStarks/lime/internal/screen"
)

type Application struct {
	Version *Version
	Screen  *screen.Screen
}

func NewApplication(v *Version) Application {
	return Application{
		Version: v,
		Screen:  screen.NewScreen(v.LimeVersion),
	}
}
