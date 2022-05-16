package editor

import (
	"bytes"
	"context"
	"os"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/screen"
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	buffer        *bytes.Buffer
	screen        *screen.Screen
	getContext   func() context.Context
	setContext func(context.Context)
	cancelContext context.CancelFunc
}

func NewEditor(screen *screen.Screen, setCtx func(context.Context), getCtx func() context.Context, cancelCtx context.CancelFunc) *Editor {
	return &Editor{
		buffer:        bytes.NewBuffer(nil),
		screen:        screen,
		getContext: getCtx,
		setContext: setCtx,
		cancelContext: cancelCtx,
	}
}

func (e *Editor) WriteBuffer(b []byte) {
	e.buffer.Write(b)
}

// b: the text to be written to the buffer;
// i: the index to insert the text at
func (e *Editor) InsertToBuffer(b []byte, i int) {
	var newBuf []byte
	if i >= len(e.buffer.Bytes()) { // If the index is greater than the length of the buffer, append the text
		newBuf = append(e.buffer.Bytes(), b...)
	} else if i < 0 { // If the index is less than 0, insert the text at the beginning of the buffer
		newBuf = append(b, e.buffer.Bytes()...)
	} else { // If the index is within the buffer, insert the text at the index
		newBuf = append(e.buffer.Bytes()[:i], append(b, e.buffer.Bytes()[i:]...)...)
	}
	e.buffer.Reset()
	e.buffer.Write(newBuf)
}

func (e *Editor) RemoveFromBuffer(i int) {
	var newBuf []byte
	if i >= len(e.buffer.Bytes()) {
		newBuf = e.buffer.Bytes()[:e.buffer.Len()]
	} else if i < 0 { // If the index is less than 0, insert the text at the beginning of the buffer
		// Do nothing
	} else {
		newBuf = append(e.buffer.Bytes()[:i], e.buffer.Bytes()[i+1:]...)
	}
	e.buffer.Reset()
	e.buffer.Write(newBuf)
}

func (e *Editor) ReadBufferByte() []byte {
	return e.buffer.Bytes()
}

func (e *Editor) ReadBufferString() string {
	return e.buffer.String()
}

func (e *Editor) CancelContext() {
	e.cancelContext()
}

func (e *Editor) Launch(file *os.File) {
	defer e.CancelContext()

	// Set initial buffer index. This is used to determine the index to write to
	ctx := context.WithValue(e.getContext(), constants.BufferIndexCtxKey, 0)
	e.setContext(ctx)

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
				lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
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
				lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
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

				lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
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

				// Read file into buffer
				e.buffer.Reset()
				_, err := e.buffer.ReadFrom(file)
				if err != nil {
					utils.LogMessage(err.Error())
				}
				e.screen.SetCursor(constants.EditorPaddingLeft+2, constants.EditorPaddingTop+1) // Set cursor to the start
				e.Read()
			} else if (ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2) && started {
				e.BackSpace()
			} else if ev.Key() == tcell.KeyEnter && started {
				e.Write('\n' ,constants.KeyArrowDown, 1)
			} else if ev.Key() == tcell.KeyTab && started {
				e.Write('\t', constants.KeyArrowRight, configs.TabSize)
			} else {
				char := ev.Rune()
				if char >= 32 && started { // Printable character
					e.Write(char, constants.KeyArrowRight, 1)
				}
			}
		}
	}
}
