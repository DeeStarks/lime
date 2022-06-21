package init

import (
	"context"
	"os"

	"github.com/DeeStarks/lime/internal/editor"
	"github.com/DeeStarks/lime/internal/screen"
)

type Application struct {
	Screen *screen.Screen
	Editor *editor.Editor
}

func NewApplication(f *os.File) Application {
	// Create a general context
	ctx, cancel := context.WithCancel(context.Background())
	setContext := func(newCtx context.Context) {
		ctx = newCtx
	}

	getContext := func() context.Context {
		return ctx
	}

	// Create a new screen
	screen := screen.NewScreen(setContext, getContext, cancel)

	return Application{
		Screen: screen,
		Editor: editor.NewEditor(f, screen, setContext, getContext, cancel),
	}
}

func (a Application) Start() {
	a.Screen.ShowBox() // Launch the initial box
	a.Editor.Launch()
}
