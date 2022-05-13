package editor

import (
	"fmt"
	"strconv"

	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) showBars(filename string) { // Display title, numbering... bars
	filename = fmt.Sprintf("•• %s ••", filename)
	sw, sh := e.screen.GetScreen().Size() // Get screen size

	// Title bar
	tstyle := utils.CreateStyle(tcell.ColorReset, tcell.ColorAntiqueWhite)
	e.screen.DrawBox(0, 0, sw, 0, "", true, false, true, tcell.StyleDefault)
	e.screen.DrawText((sw/2)-(len(filename)/2)-1, 0, (sw/2)+(len(filename)/2)+1, 0, filename, tstyle)

	// Draw line numbering side bar
	nw, nh := 6, sh
	nstyle := utils.CreateStyle(tcell.ColorReset, tcell.ColorDimGray)
	e.screen.DrawBox(0, 1, nw, nh, "", true, false, false, tcell.StyleDefault)
	count := sh
	for i := 1; i <= count; i++ {
		e.screen.DrawText(2, i, nw, i+1, strconv.Itoa(i), nstyle)
	}
}
