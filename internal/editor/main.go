package editor

import (
	"context"
	"os"

	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/screen"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen        *screen.Screen
	context       context.Context
	cancelContext context.CancelFunc
}

func NewEditor(screen *screen.Screen) *Editor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Editor{
		screen:        screen,
		context:       ctx,
		cancelContext: cancel,
	}
}

func (e *Editor) GetContext() context.Context {
	return e.context
}

func (e *Editor) CancelContext() {
	e.cancelContext()
}

func (e *Editor) Launch(file *os.File) {
	ctx := e.GetContext()
	defer e.CancelContext()

	// Set initial line counter to 0
	e.context = context.WithValue(ctx, constants.LineCounterCtxKey, NewLineCounter(0))

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
			} else if ev.Key() == tcell.KeyLeft && started { // Left
				e.screen.MoveCursor(constants.KeyArrowLeft)
			} else if ev.Key() == tcell.KeyRight && started { // Right
				// Make sure the cursor does not exceed the last character
				lc := e.context.Value(constants.LineCounterCtxKey).(LineCounter)
				x, y := e.screen.GetCursorPosition()
				// Subtract paddings and extra paddings
				x = x - (constants.EditorPaddingLeft + 2)
				y = y - (constants.EditorPaddingTop + 1)

				if len(lc) > y && x < lc[y] {
					e.screen.MoveCursor(constants.KeyArrowRight)
				}
			} else if ev.Key() == tcell.KeyUp && started { // Up
				e.screen.MoveCursor(constants.KeyArrowUp)

				// To ensure cursor position doens't exceed the last character
				lc := e.context.Value(constants.LineCounterCtxKey).(LineCounter)
				x, y := e.screen.GetCursorPosition()
				x = x - (constants.EditorPaddingLeft + 2)
				y = y - (constants.EditorPaddingTop + 1)

				if lc[y] < x {
					distance := x - lc[y]
					for i := 0; i < distance; i++ {
						e.screen.MoveCursor(constants.KeyArrowLeft)
					}
				}
			} else if ev.Key() == tcell.KeyDown && started { // Down
				_, y := e.screen.GetCursorPosition()
				// Subtract paddings and extra paddings
				y = y - (constants.EditorPaddingTop + 1)

				lc := e.context.Value(constants.LineCounterCtxKey).(LineCounter)
				// Make sure the cursor does not exceed the last line
				if y < len(lc) {
					e.screen.MoveCursor(constants.KeyArrowDown)
				}

				// Cursor position doens't exceed the last character
				// First get new cursor position
				x, y := e.screen.GetCursorPosition()
				x = x - (constants.EditorPaddingLeft + 2)
				y = y - (constants.EditorPaddingTop + 1)

				if len(lc) > y && lc[y] < x {
					distance := x - lc[y]
					for i := 0; i < distance; i++ {
						e.screen.MoveCursor(constants.KeyArrowLeft)
					}
				}
			} else if ev.Key() == tcell.KeyCtrlW {
				e.screen.GetScreen().Sync()
				if !started {
					e.screen.GetScreen().Clear()
					e.showBars(file.Name())
					started = true
				}

				e.ReadFile(file)
			}
		}
	}
}
