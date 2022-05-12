package editor

import (
	"github.com/DeeStarks/lime/internal/screen"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen 	*screen.Screen
}

func NewEditor(screen *screen.Screen) *Editor {
	return &Editor{
		screen,
	}
}

func (e *Editor) Launch() {
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
				e.screen.GetScreen().Clear()

				// TODO: Implement editor
				
			}
		}
	}
}