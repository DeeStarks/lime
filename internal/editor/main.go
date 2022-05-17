package editor

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/screen"
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	file          *os.File
	initialBuffer *bytes.Buffer // This will be used to know if there are changes made, hence "modified"
	buffer        *bytes.Buffer
	lines         [][][]byte // Each line on the screen
	screen        *screen.Screen
	width         int
	height        int
	title         string
	lineNumbering []LineNumbering
	defStyle      tcell.Style
	getContext    func() context.Context
	setContext    func(context.Context)
	cancelContext context.CancelFunc
}

func NewEditor(file *os.File, screen *screen.Screen, setCtx func(context.Context), getCtx func() context.Context, cancelCtx context.CancelFunc) *Editor {
	sw, sh := screen.GetScreen().Size()
	// Make initial numbering to fill page
	nl := make([]LineNumbering, sh)
	for i := 0; i < sh; i++ {
		nl[i] = LineNumbering{
			number: i + 1,
			yAxis:  i + (constants.EditorPaddingTop + 1),
		}
	}

	// Read the file into a buffer
	b, err := ioutil.ReadFile(file.Name())
	if err != nil {
		utils.LogMessage("Error reading file: " + err.Error())
	}
	return &Editor{
		file:          file,
		initialBuffer: bytes.NewBuffer(b),
		buffer:        bytes.NewBuffer(b),
		lines:         [][][]byte{},
		screen:        screen,
		width:         sw - constants.EditorPaddingLeft - constants.EditorPaddingRight,
		height:        sh - constants.EditorPaddingTop - constants.EditorPaddingBottom,
		title:         file.Name(),
		lineNumbering: nl,
		defStyle:      utils.CreateStyle(tcell.ColorBlack, tcell.ColorWhiteSmoke),
		getContext:    getCtx,
		setContext:    setCtx,
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

func (e *Editor) ReadInitialBufferByte() []byte {
	return e.initialBuffer.Bytes()
}

func (e *Editor) ReadBufferByte() []byte {
	return e.buffer.Bytes()
}

func (e *Editor) ReadInitialBufferString() string {
	return e.initialBuffer.String()
}

func (e *Editor) ReadBufferString() string {
	return e.buffer.String()
}

// This should be called on every cursor movement.
// It helps to hold the index where new input will be appended
func (e *Editor) UpdateBufferIndex() {
	// Spread the lines into a new array.
	// This will hold the line number and the line itself
	// This is used because of the wrapping of lines
	var newLines [][][]byte
	for i, line := range e.lines {
		for _, l := range line {
			// "i" is the line number
			// It will be stored as a byte type first
			iByte := []byte(strconv.Itoa(i))
			// Append the line number to the line
			newLines = append(newLines, [][]byte{iByte, l})
		}
	}

	// lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
	x, y := e.screen.GetCursorPosition()

	// Get the line number of the cursor
	prevs := e.getContext().Value(constants.StartLineCtxKey).(int) + y // The previous lines, including the ones scrolled off the screen
	var count int
	if len(newLines) > 0 { // If there are lines in the buffer
		for i := 0; i < prevs-1; i++ {
			currLine, nextLine := newLines[i][0][0], newLines[i+1][0][0]
			if currLine == nextLine {
				count += len(newLines[i][1])
			} else {
				count += len(newLines[i][1]) + 1 // +1 for the newline
				// currentLine = newLines[i][0][0]
			}
		}
	}
	count += x // The current line
	count -= constants.EditorPaddingLeft + 2

	// Update the line counter
	ctx := context.WithValue(e.getContext(), constants.BufferIndexCtxKey, count-1)
	e.setContext(ctx)
}

func (e *Editor) CancelContext() {
	e.cancelContext()
}

func (e *Editor) Launch() {
	defer e.CancelContext()

	// Set Initial contexts
	func() {
		// Set initial buffer index. This is used to determine the index to write to
		ctx := context.WithValue(e.getContext(), constants.BufferIndexCtxKey, 0)
		e.setContext(ctx)

		// Set initial line counter
		ctx = context.WithValue(e.getContext(), constants.LineCounterCtxKey, LineCounter{})
		e.setContext(ctx)

		// Set starting line context
		ctx = context.WithValue(e.getContext(), constants.StartLineCtxKey, 0)
		e.setContext(ctx)

		// Set total lines context
		ctx = context.WithValue(e.getContext(), constants.TotalLinesCtxKey, 0)
		e.setContext(ctx)
	}()

	var started bool
	for {
		e.screen.GetScreen().Show()
		event := e.screen.GetScreen().PollEvent()

		switch ev := event.(type) {
		case *tcell.EventResize:
			if !started {
				e.screen.ShowBox()
			} else {
				// Update screen size
				sw, sh := e.screen.GetScreen().Size()
				e.width = sw - constants.EditorPaddingLeft - constants.EditorPaddingRight
				e.height = sh - constants.EditorPaddingTop - constants.EditorPaddingBottom

				// Sync
				e.showBars()
				e.Read()
			}
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

				if len(lc) > y {
					if lc[y] < x {
						distance := x - lc[y]
						for i := 0; i < distance; i++ {
							e.screen.MoveCursor(constants.KeyArrowLeft)
						}
					}
				} else {
					e.screen.SetCursor(constants.EditorPaddingLeft+2, y-1)
				}
				e.Read() // Update incase page has been moved up
			} else if ev.Key() == tcell.KeyDown && started { // Down
				// _, y := e.screen.GetCursorPosition()
				// Subtract paddings and extra paddings
				// y = y - (constants.EditorPaddingTop + 1)

				lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
				// Make sure the cursor does not exceed the last line
				e.screen.MoveCursor(constants.KeyArrowDown)
				// if y < len(lc)-1 {
				// }

				// Cursor position doens't exceed the last character
				// First get new cursor position
				x, y := e.screen.GetCursorPosition()
				x = x - (constants.EditorPaddingLeft + 2)
				y = y - (constants.EditorPaddingTop + 1)

				if len(lc)-1 > y && lc[y] < x {
					distance := x - lc[y]
					for i := 0; i < distance; i++ {
						e.screen.MoveCursor(constants.KeyArrowLeft)
					}
				}
				e.Read() // Update incase page has been moved down
			} else if ev.Key() == tcell.KeyCtrlS && started {
				e.Save()
			} else if ev.Key() == tcell.KeyCtrlW {
				e.screen.GetScreen().Sync()
				if !started {
					e.screen.GetScreen().Clear()
					e.showBars()
					started = true
				}

				// Set cursor to the start
				e.screen.SetCursor(constants.EditorPaddingLeft+2, constants.EditorPaddingTop+1)

				e.Read()
			} else if (ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2) && started {
				e.BackSpace()
			} else if ev.Key() == tcell.KeyEnter && started {
				e.Write('\n')
			} else if ev.Key() == tcell.KeyTab && started {
				e.Write('\t')
			} else {
				char := ev.Rune()
				if char >= 32 && started { // Printable character
					e.Write(char)
				}
			}
		}
		e.UpdateBufferIndex() // Always update the buffer index
	}
}
