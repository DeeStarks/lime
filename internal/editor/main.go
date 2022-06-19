package editor

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/highlighters"
	"github.com/DeeStarks/lime/internal/screen"
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	file          *os.File
	initialBuffer *bytes.Buffer
	buffer        *bytes.Buffer
	insIndex      int   // The index where new input will be appended
	startLine     int   // The line first line on the screen
	lineCounter   []int // The length of each line
	screen        *screen.Screen
	width         int
	height        int
	title         string
	lineNumbering []LineNumbering
	defStyle      tcell.Style
	getContext    func() context.Context
	setContext    func(context.Context)
	cancelContext context.CancelFunc
	highlighters  *highlighters.Highlighter
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
		insIndex:      0,
		startLine:     0,
		lineCounter:   make([]int, sh),
		screen:        screen,
		width:         sw - constants.EditorPaddingLeft - constants.EditorPaddingRight,
		height:        sh - constants.EditorPaddingTop - constants.EditorPaddingBottom,
		title:         file.Name(),
		lineNumbering: nl,
		defStyle:      utils.CreateStyle(tcell.ColorBlack, tcell.ColorWhiteSmoke),
		getContext:    getCtx,
		setContext:    setCtx,
		cancelContext: cancelCtx,
		highlighters:  highlighters.NewHighlighter(utils.GetFileExtension(file.Name())),
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

func (e *Editor) UpdateCursorPosition() {
	var (
		x        = constants.EditorPaddingLeft + 2
		y        = constants.EditorPaddingTop + 1
		currLine int
	)

	for i := 0; i < e.insIndex; i++ {
		if e.buffer.Bytes()[i] == '\n' || x-(constants.EditorPaddingLeft+2) > e.lineCounter[y-(constants.EditorPaddingTop+1)] {
			currLine++
			x = constants.EditorPaddingLeft + 1

			if currLine > e.startLine {
				y++
			}
		}

		if currLine >= e.startLine {
			if e.ReadBufferByte()[i] == '\t' {
				x += configs.TabSize
			} else {
				x++
			}
		}
	}

	e.screen.SetCursor(x, y)
}

func (e *Editor) CancelContext() {
	e.cancelContext()
}

func (e *Editor) Launch() {
	defer e.CancelContext()

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
			} else if ev.Key() == tcell.KeyLeft && started && e.insIndex > 0 { // Left
				e.insIndex--
				e.UpdateCursorPosition()
			} else if ev.Key() == tcell.KeyRight && started && e.insIndex < len(e.ReadBufferByte()) { // Right
				e.insIndex++
				e.UpdateCursorPosition()
			} else if ev.Key() == tcell.KeyUp && started { // Up
				// Get the width between the index and beginning of the line
				x, _ := e.screen.GetCursorPosition()
				width := x - (constants.EditorPaddingLeft - 2)

				// Count two newlines backwards and add the width
				var (
					count    int
					prevLine int
				)
				for i := e.insIndex - 1; i >= 0; i-- {
					if e.ReadBufferByte()[i] == '\n' {
						if count == 1 {
							prevLine = i
							break
						}
						count++
					}
				}
				e.insIndex = prevLine + width
				e.UpdateCursorPosition()
			} else if ev.Key() == tcell.KeyDown && started { // Down
				// Get the width between the index and beginning of the line
				x, _ := e.screen.GetCursorPosition()
				width := x - (constants.EditorPaddingLeft - 2)

				// Get index of next line
				var nextLine int
				for i := e.insIndex; i < len(e.ReadBufferByte()); i++ {
					if e.ReadBufferByte()[i] == '\n' {
						nextLine = i + 1
						break
					}
				}
				e.insIndex = nextLine + width
				e.UpdateCursorPosition()
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
	}
}
