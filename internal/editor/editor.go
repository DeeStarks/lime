package editor

import (
	"io/ioutil"
	"os"

	"github.com/DeeStarks/lime/configs"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) Edit(file *os.File) {
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

	var line int
	var lineCount int
	for _, c := range content {
		switch c {
		case '\n': // Newline
			line++
			lineCount = 0 // Start from the beginning of the line
		case '\t': // Tab
			lineCount += configs.TabSize
		default:
			lineCount++
		}
		screen.SetContent(coord.x+lineCount, coord.y+line, rune(c), nil, tcell.StyleDefault)
	}
	screen.Show()
	screen.Sync()
}

func drawEditBox(e *Editor) (tcell.Screen, struct{x, y int}) {
	sw, sh := e.screen.GetScreen().Size()
	coord := struct {
		x, y int
	} {6, 0}
	e.screen.DrawBox(coord.x, coord.y, sw-1, sh, "", false, false, true, tcell.StyleDefault)
	return e.screen.GetScreen(), coord
}