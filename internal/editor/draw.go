package editor

import (
	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func drawEditBox(e *Editor) (tcell.Screen, struct{ x, y int }) {
	sw, sh := e.screen.GetScreen().Size()
	coord := struct {
		x, y int
	}{constants.EditorPaddingLeft, constants.EditorPaddingTop}
	e.screen.DrawBox(coord.x, coord.y, sw-1, sh, "", false, false, false, utils.CreateStyle(tcell.ColorYellow, tcell.ColorWheat))
	return e.screen.GetScreen(), coord
}
