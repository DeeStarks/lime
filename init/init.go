package init

import (
	"context"
	"os"

	"github.com/DeeStarks/lime/internal/editor"
	"github.com/DeeStarks/lime/internal/screen"
)

type Application struct {
	Version *Version
	Screen  *screen.Screen
	Editor  *editor.Editor
}

func NewApplication(f *os.File, v *Version) Application {
	// Create a general context
	ctx, cancel := context.WithCancel(context.Background())
	setContext := func(newCtx context.Context) {
		ctx = newCtx
	}

	getContext := func() context.Context {
		return ctx
	}

	// Create a new screen
	screen := screen.NewScreen(v.LimeVersion, setContext, getContext, cancel)

	return Application{
		Version: v,
		Screen:  screen,
		Editor:  editor.NewEditor(f, screen, setContext, getContext, cancel),
	}
}

func (a Application) Start() {
	a.Screen.ShowBox() // Launch the initial box
	a.Editor.Launch()
}
