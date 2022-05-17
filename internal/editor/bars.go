package editor

import (
	"strconv"

	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type LineNumbering struct {
	number int
	yAxis  int
}

func (e *Editor) showBars() { // Display title, numbering... bars
	sw, sh := e.screen.GetScreen().Size() // Get screen size
	style := utils.CreateStyle(tcell.ColorReset, tcell.ColorDarkGray)

	// Title bar
	e.screen.DrawBox(0, 0, sw, 0, "", true, false, false, tcell.StyleDefault)
	e.screen.DrawText((sw/2)-(len(e.title)/2)-1, 0, (sw/2)+(len(e.title)/2)+1, 0, e.title, style)

	// Draw line numbering side bar by
	nw, nh := 6, sh
	e.screen.DrawBox(0, 1, nw, nh, "", true, false, false, tcell.StyleDefault)
	for _, l := range e.lineNumbering {
		e.screen.DrawText(1, l.yAxis, nw, l.yAxis+1, strconv.Itoa(l.number), style)
	}

	// Fill the remain space with a dummy character
	prev := e.lineNumbering[len(e.lineNumbering)-1].yAxis + 1
	style = utils.CreateStyle(tcell.ColorReset, tcell.ColorDimGray)
	for i := prev; i < sh; i++ {
		e.screen.DrawText(1, i, nw, i+1, "~", style)
	}

}
