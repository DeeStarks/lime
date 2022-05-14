package editor

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/constants"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) ReadFile(file *os.File) {
	screen, coord := drawEditBox(e)
	// Add padding
	coord.y += 1
	coord.x += 1
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		screen.Fill(' ', tcell.StyleDefault)
		screen.Show()
		return
	}

	// Place the cursor at the beginning of the file
	e.screen.SetCursor(coord.x+1, coord.y)

	// Draw the content
	var line int
	var lineCount int
	for _, c := range content {
		switch c {
		case '\n': // Newline
			line++
			lineCount = 0 // Start from the beginning of the line

			// Update the line counter
			lc := e.context.Value(constants.LineCounterCtxKey).(LineCounter)
			newCtx := append(lc, lineCount)
			e.context = context.WithValue(e.context, constants.LineCounterCtxKey, newCtx)
		case '\t': // Tab
			lineCount += configs.TabSize
			// Update the current line counter by the tab size
			lc := e.context.Value(constants.LineCounterCtxKey).(LineCounter)
			lc[line] = lc[line] + configs.TabSize
			e.context = context.WithValue(e.context, constants.LineCounterCtxKey, lc)
		default:
			lineCount++
			// Update the counter for the current line
			lc := e.context.Value(constants.LineCounterCtxKey).(LineCounter)
			lc[line] = lc[line] + 1
			e.context = context.WithValue(e.context, constants.LineCounterCtxKey, lc)
		}
		screen.SetContent(coord.x+lineCount, coord.y+line, rune(c), nil, tcell.StyleDefault)
	}
	screen.Show()
	screen.Sync()
}

func drawEditBox(e *Editor) (tcell.Screen, struct{ x, y int }) {
	sw, sh := e.screen.GetScreen().Size()
	coord := struct {
		x, y int
	}{constants.EditorPaddingLeft, constants.EditorPaddingTop}
	e.screen.DrawBox(coord.x, coord.y, sw-1, sh, "", false, false, false, tcell.StyleDefault)
	return e.screen.GetScreen(), coord
}

func (e *Editor) WriteFile() {

}
