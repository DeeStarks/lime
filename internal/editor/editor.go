package editor

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/utils"
)

func (e *Editor) Read() {
	// Clear the editor
	e.Clear()

	screen, _ := drawEditBox(e)
	// We'll create a 3D array to store the bytes of each line
	// The first dimension is each line
	// The second dimension is another array of lines for line wrapping
	// The third dimension is the bytes of each line
	var lines [][][]byte
	buf := e.ReadBufferString()
	// Replace tabs with spaces
	spaces := make([]byte, configs.TabSize)
	for i := 0; i < configs.TabSize; i++ {
		spaces[i] = ' '
	}
	buf = strings.ReplaceAll(buf, "\t", string(spaces))

	for _, l := range strings.Split(buf, "\n") {
		cl := append([][]byte{}, []byte(l))
		lines = append(lines, cl)
	}

	// We'll now wrap the lines
	var totalLineLength int
	for i, l := range lines {
		maxLineLength := e.width - 2
		totalLineLength++
		if len(l[0]) > maxLineLength {
			// We'll break the line into length of the width
			// and store it in the lines array
			var wrappedLines [][]byte
			left := 0
			right := maxLineLength

			for {
				if right >= len(l[0]) {
					wrappedLines = append(wrappedLines, l[0][left:len(l[0])])
					break
				}

				wrappedLines = append(wrappedLines, l[0][left:right])
				left = right
				right += maxLineLength
				totalLineLength++
			}

			lines[i] = wrappedLines
		}
	}
	utils.LogMessage("TLC: %d", totalLineLength)

	// Save the lines in the editor
	e.lines = lines

	var numbering []LineNumbering

	// We'll now draw the lines
	startLine := e.getContext().Value(constants.StartLineCtxKey).(int)
	numberCount := startLine
	currLine := 0
	yAxis := constants.EditorPaddingTop + 1
	endLine := e.height - 1
	if len(lines) < endLine {
		endLine = len(lines)
	}

	// Restart the line counter
	e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, LineCounter{}))
	for _, l := range lines[startLine:endLine] {
		// Create a new line counter
		nmb := struct {
			number int
			yAxis  int
		}{
			number: numberCount + 1,
			yAxis:  yAxis,
		}

		for j, wl := range l {
			// Increment yaxis if j more than one line
			if j > 0 {
				yAxis++
				currLine++
			}

			// Draw the line
			e.screen.DrawText(
				constants.EditorPaddingLeft+2,
				currLine+(constants.EditorPaddingTop+1),
				e.width+50,
				e.height, string(wl), e.defStyle)

			// Increment the line counter
			lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
			newCtr := append(lc, len(wl))
			e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, newCtr))
		}
		yAxis++
		currLine++
		numberCount++
		// Append numbering
		numbering = append(numbering, nmb)
	}

	// Update bars
	e.lineNumbering = numbering

	if e.ReadBufferString() != e.ReadInitialBufferString() {
		e.title = fmt.Sprintf("•• %s (modified) ••", e.file.Name())
	} else {
		e.title = fmt.Sprintf("•• %s ••", e.file.Name())
	}

	screen.Show()
	e.showBars()
}

func (e *Editor) Clear() {
	sw, sh := e.screen.GetScreen().Size()
	ew, eh := sw-constants.EditorPaddingLeft-constants.EditorPaddingRight, sh-constants.EditorPaddingTop-constants.EditorPaddingBottom

	// Create a new buffer filled with spaces and draw on the editor
	buf := make([]byte, ew)
	for i := 0; i < ew; i++ {
		buf[i] = ' '
	}

	// Draw the content
	for line := constants.EditorPaddingTop + 1; line <= eh; line++ {
		e.screen.DrawText(constants.EditorPaddingLeft, line, sw-constants.EditorPaddingRight, line+1, string(buf), e.defStyle)
	}
}

// cd: Cursor direction -> left or right or up or down
//
// mt: Number of times to move the cursor
func (e *Editor) Write(char rune) {
	cx, cy := e.screen.GetCursorPosition()

	// Write to buffer
	currIndex := e.getContext().Value(constants.BufferIndexCtxKey).(int) + 1

	switch char {
	case '\t':
		e.screen.SetCursor(cx+configs.TabSize, cy)
		// We use spaces instead of tabs due BufferIndex counting number of configs.TabSize
		// instead of the length of a tab character '\t' = 1. This may be changed in the future
		tab := make([]byte, configs.TabSize)
		for i := 0; i < configs.TabSize; i++ {
			tab[i] = ' '
		}
		e.InsertToBuffer(tab, currIndex)
	case '\n':
		e.screen.SetCursor(constants.EditorPaddingLeft+2, cy+1)
		e.InsertToBuffer([]byte(string(char)), currIndex)
	default:
		e.screen.MoveCursor(constants.KeyArrowRight)
		e.InsertToBuffer([]byte(string(char)), currIndex)
	}

	// Move cursor down if we're at the end of the line
	cx, cy = e.screen.GetCursorPosition()
	sw, _ := e.screen.GetScreen().Size()
	
	if cx == sw-constants.EditorPaddingRight-1 {
		e.screen.SetCursor(constants.EditorPaddingLeft+2, cy+1)
		// currIndex = e.getContext().Value(constants.BufferIndexCtxKey).(int) + 1
		utils.LogMessage("Char: %s", string(e.ReadBufferByte()))
	}
	e.Read()
}

func (e *Editor) BackSpace() {
	currIndex := e.getContext().Value(constants.BufferIndexCtxKey).(int)
	if currIndex < 0 {
		return
	}

	e.RemoveFromBuffer(currIndex)

	// Update the line counter
	cx, cy := e.screen.GetCursorPosition()
	cy = cy - (constants.EditorPaddingTop + 1) // Remove padding to get the real line number inside the editor
	lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
	lc[cy]--
	e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, lc))

	// Read to editor
	e.Read()

	// Move cursor
	if cx <= constants.EditorPaddingLeft+2 {
		e.screen.SetCursor(lc[cy-1]+constants.EditorPaddingLeft+2, cy+constants.EditorPaddingTop)
	} else {
		e.screen.MoveCursor(constants.KeyArrowLeft)
	}
}

func (e *Editor) Save() {
	bt := e.ReadBufferByte()
	e.initialBuffer = bytes.NewBuffer(bt)
	ioutil.WriteFile(e.file.Name(), e.ReadInitialBufferByte(), 0644)
	e.Read() // To update title bar
}