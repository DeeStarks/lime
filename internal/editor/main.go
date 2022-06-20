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
	buffer        bytes.Buffer
	history []bytes.Buffer
	maxHistory    int
	currentBuffer int 	// The current buffer in the history
	insIndex      int   // The index where new input will be appended
	indexHistory  []int // The index history
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
	// Initialize histories
	maxHistory := 50 // The maximum number of history buffers
	history := make([]bytes.Buffer, maxHistory)
	history[0].Write(bytes.NewBuffer(b).Bytes())

	indexHistory := make([]int, maxHistory)
	indexHistory[0] = 0

	return &Editor{
		file:          file,
		initialBuffer: bytes.NewBuffer(b),
		history: 		history,
		maxHistory:    maxHistory,
		buffer:        history[0],
		currentBuffer: 0,
		insIndex:      0,
		indexHistory:  indexHistory,
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

// b: the text to be written to the buffer;
// i: the index to insert the text at
func (e *Editor) InsertToBuffer(b []byte, i int) {
	if i >= len(e.ReadBufferByte()) { // If the index is greater than the length of the buffer, append the text
		e.buffer.Write(b)
		utils.LogMessage("New Buff: %v", e.buffer.String())
	} else if i < 0 { // If the index is less than 0, insert the text at the beginning of the buffer
		e.buffer.Reset()
		e.buffer.Write(append(b, e.ReadBufferByte()...))
	} else { // If the index is within the buffer, insert the text at the index
		prevBuf := e.ReadBufferByte()
		e.buffer.Reset()
		e.buffer = *bytes.NewBuffer(append(prevBuf[:i], append(b, prevBuf[i:]...)...))
	}

	// Update the buffer history
	if e.currentBuffer + 1 >= e.maxHistory { // Make sure the history doesn't grow too large
		// If the buffer history is full, remove the oldest buffer
		newHistory := make([]bytes.Buffer, e.maxHistory)
		copy(newHistory, e.history[1:])
		newHistory[e.maxHistory-1].Reset()
		newHistory[e.maxHistory-1].Write(e.buffer.Bytes())
		e.history = newHistory

		// Update the index history
		newIndexHistory := make([]int, e.maxHistory)
		copy(newIndexHistory, e.indexHistory[1:])
		newIndexHistory[e.maxHistory-1] = e.insIndex
		e.indexHistory = newIndexHistory
	} else {
		// If the buffer history is not full, add a new buffer
		e.currentBuffer++
		e.history[e.currentBuffer].Reset()
		e.history[e.currentBuffer].Write(e.ReadBufferByte())
		// Update all history buffers after the current buffer to null
		if e.currentBuffer + 1 < e.maxHistory {
			for i := e.currentBuffer + 1; i < len(e.history); i++ {
				e.history[i] = *bytes.NewBuffer(nil)
			}
		}

		// Update the index history
		e.indexHistory[e.currentBuffer] = e.insIndex
		// Update the index history after the current buffer to 0
		if e.currentBuffer + 1 < e.maxHistory {
			for i := e.currentBuffer + 1; i < len(e.indexHistory); i++ {
				e.indexHistory[i] = 0
			}
		}
	}

}

func (e *Editor) Undo() {
	prev := e.currentBuffer - 1
	if prev >= 0 {
		e.currentBuffer = prev
		e.buffer.Reset()
		e.buffer.Write(e.history[prev].Bytes())

		// Decrement the index
		e.insIndex = e.indexHistory[prev] + 1
		e.Read()
		e.UpdateCursorPosition()
	}
}

func (e *Editor) Redo() {
	next := e.currentBuffer + 1
	if next < len(e.history) && e.history[next].Len() > 0 {
		e.currentBuffer = next
		e.buffer.Reset()
		e.buffer.Write(e.history[next].Bytes())

		// Increment the index
		e.insIndex = e.indexHistory[next] + 1
		e.Read()
		e.UpdateCursorPosition()
	}
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

func (e *Editor) ScrollUp() {
	// This checks if that the cursor is at the last line before scrolling
	_, y := e.screen.GetCursorPosition()
	_, sh := e.screen.GetScreen().Size()
	sh = sh - (constants.EditorPaddingTop+1)
	line := y - (constants.EditorPaddingTop+1)

	if line >= sh && e.insIndex <= len(e.ReadBufferByte()) {
		e.startLine++
		e.Read()
		e.UpdateCursorPosition()
	}
}

func (e *Editor) ScrollDown() {
	// This checks if that the cursor is at the first line before scrolling
	_, y := e.screen.GetCursorPosition()
	line := y - (constants.EditorPaddingTop + 1)

	if line == 0 && e.startLine > 0 {
		e.startLine--
		e.Read()
		e.UpdateCursorPosition()
	}
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
				e.UpdateCursorPosition()
				e.Read()
			}
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyCtrlQ {
				e.screen.Quit()
			} else if ev.Key() == tcell.KeyLeft && started && e.insIndex > 0 { // Left
				e.insIndex--
				e.UpdateCursorPosition()
				e.ScrollDown()
			} else if ev.Key() == tcell.KeyRight && started && e.insIndex < len(e.ReadBufferByte()) { // Right
				e.insIndex++
				e.UpdateCursorPosition()
				e.ScrollUp()
			} else if ev.Key() == tcell.KeyUp && started { // Up
				var width int

				// Get the width from the beginning of the line to the index
				x, y := e.screen.GetCursorPosition()
				y -= constants.EditorPaddingTop + 1
				width = x - (constants.EditorPaddingLeft + 2)

				if e.startLine > 0 && y <= 0 {
					if y - 1 < 0 {
						e.startLine--
						e.Read()
					}
				}

				if y > 0 {
					// First, we remove the width from the beginning of the current line
					prevIndex := e.insIndex
					e.insIndex -= width
					if e.insIndex < 0 {
						e.insIndex = 0
					}

					// Count number of tabs between the previous and current index
					var tabs int
					for i := prevIndex; i > e.insIndex; i-- {
						if e.ReadBufferByte()[i] == '\t' {
							tabs++
						}
					}

					// We incrment the insIndex by the number of tabs
					if tabs > 0 {
						e.insIndex += tabs * configs.TabSize - 1
					}

					if e.insIndex > 0 {
						e.insIndex--

						width = e.lineCounter[y-1] - width
						if width <= 0 {
							width = 0
						}
						prevIndex = e.insIndex
						e.insIndex -= width
						
						if e.insIndex < 0 { 
							// Just a precaution to make sure insertion index doesn't go below 0
							e.insIndex = 0
						}

						// // Count tabs between the previous and current index
						tabs = 0
						for i := e.insIndex; i < prevIndex; i++ {
							utils.LogMessage("%s", string(e.ReadBufferByte()[i]))
							if e.ReadBufferByte()[i] == '\t' {
								tabs++
							}
						}

						// // We incrment the insIndex by the number of tabs
						if tabs > 0 {
							e.insIndex += tabs * configs.TabSize - 1
						}
					}
				}

				// // Update cursor position again to make sure it's in the right place
				e.UpdateCursorPosition()
			} else if ev.Key() == tcell.KeyDown && started { // Down
				var (
					width int // Width from the beginning of the line to the index
					rem  int // Remaining width to be added to the index
				)

				for i := e.insIndex; i < len(e.ReadBufferByte()); i++ {
					if e.ReadBufferByte()[i] == '\n' {
						break
					}
					rem++
				}

				// Get the width between the index and beginning of the line
				_, sh := e.screen.GetScreen().Size()
				x, y := e.screen.GetCursorPosition()
				width = x - (constants.EditorPaddingLeft + 2)

				// Check if next line is longer than the current one before moving to the next line
				if y < (sh - 1) && y < len(e.lineCounter) {
					var nextLine int

					if e.lineCounter[y] >= width {
						nextLine = e.insIndex + rem + width + 1 // +1 for newline
					} else {
						nextLine = e.insIndex + rem + e.lineCounter[y] + 1 // +1 for newline
					}

					if nextLine < len(e.ReadBufferByte()) {
						e.insIndex = nextLine
					} else {
						e.insIndex = len(e.ReadBufferByte())
					}
				} else {
					if e.insIndex + rem + width + 1 < len(e.ReadBufferByte()) {
						e.insIndex += rem + width + 1 // +1 for newline
					} else {
						e.insIndex = len(e.ReadBufferByte())
					}

					// Scroll
					_, y = e.screen.GetCursorPosition()
					if y >= (sh - 1) {
						e.startLine++
						e.Read()
					}
					
				}

				// // Update cursor position again to make sure it's in the right place
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
			} else if ev.Key() == tcell.KeyCtrlZ && started {
				e.Undo()
			} else if ev.Key() == tcell.KeyCtrlY && started {
				e.Redo()
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
