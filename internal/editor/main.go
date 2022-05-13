package editor

import (
	"os"

	"github.com/DeeStarks/lime/internal/screen"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen *screen.Screen
}

func NewEditor(screen *screen.Screen) *Editor {
	return &Editor{
		screen,
	}
}

func (e *Editor) Launch(file *os.File) {
	var started bool
	for {
		e.screen.GetScreen().Show()
		event := e.screen.GetScreen().PollEvent()

		switch ev := event.(type) {
		case *tcell.EventResize:
			e.screen.GetScreen().Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyCtrlQ {
				e.screen.Quit()
			} else {
				e.screen.GetScreen().Sync()
				if !started {
					e.screen.GetScreen().Clear()
					e.showBars(file.Name())
					started = true
				}

				e.Edit(file)
			}
		}
	}
}
