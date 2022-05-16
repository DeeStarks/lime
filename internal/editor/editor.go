package editor

import (
	"context"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) Read() {
	// Clear the editor
	e.Clear()

	screen, coord := drawEditBox(e)
	// Add padding
	coord.y += 1
	coord.x += 1

	// Draw the content
	var line int
	var lineCount int
	
	// Reset the line counter
	ctx := context.WithValue(e.getContext(), constants.LineCounterCtxKey, NewLineCounter(0))
	e.setContext(ctx)

	for _, c := range e.ReadBufferString() {
		switch c {
		case '\n': // Newline
			line++
			lineCount = 0 // Start from the beginning of the line

			// Update the line counter
			lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
			newCtr := append(lc, lineCount)
			e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, newCtr))
		case '\t': // Tab
			lineCount += configs.TabSize
			// Update the current line counter by the tab size
			lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
			lc[line] = lc[line] + configs.TabSize
			e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, lc))
		default:
			lineCount++
			// Update the counter for the current line
			lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
			lc[line] = lc[line] + 1
			e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, lc))
		}
		screen.SetContent(coord.x+lineCount, coord.y+line, rune(c), nil, tcell.StyleDefault)
	}
	screen.Show()
	screen.Sync()
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
	for line := constants.EditorPaddingTop+1; line <= eh; line++ {
		e.screen.DrawText(constants.EditorPaddingLeft, line, sw-constants.EditorPaddingRight, line+1, string(buf), tcell.StyleDefault)
	}
}

// cd: Cursor direction -> left or right or up or down
// 
// mt: Number of times to move the cursor
func (e *Editor) Write(char rune, cd constants.KeyCode, mt int) {
	cx, cy := e.screen.GetCursorPosition()

	// Write to buffer
	currIndex := e.getContext().Value(constants.BufferIndexCtxKey).(int)
	e.InsertToBuffer([]byte(string(char)), currIndex)

	// Write the character to the screen
	e.screen.DrawText(cx, cy, 0, 0, string(char), tcell.StyleDefault)

	// Add to line counter
	cy = cy - (constants.EditorPaddingTop + 1) // Remove padding to get the real line number inside the editor
	lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
	switch cd {
	case constants.KeyArrowDown:
		lc = append(lc[:cy], append([]int{0}, lc[cy+1:]...)...)
		e.screen.SetCursor(constants.EditorPaddingLeft+2, cy+1)
	default:
		lc[cy] += mt
		for i := 0; i < mt; i++ {
			e.screen.MoveCursor(cd)
		}
	}
	e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, lc))

	e.Read()
}

func (e *Editor) BackSpace() {
	currIndex := e.getContext().Value(constants.BufferIndexCtxKey).(int)
	utils.LogMessage("Backspace: %d", currIndex)
	if currIndex < 0 {
		return
	}

	e.RemoveFromBuffer(currIndex)

	// Update the line counter
	_, cy := e.screen.GetCursorPosition()
	cy = cy - (constants.EditorPaddingTop + 1) // Remove padding to get the real line number inside the editor
	lc := e.getContext().Value(constants.LineCounterCtxKey).(LineCounter)
	lc[cy]--
	e.setContext(context.WithValue(e.getContext(), constants.LineCounterCtxKey, lc))

	// Read to editor
	e.Read()
	e.screen.MoveCursor(constants.KeyArrowLeft)
}