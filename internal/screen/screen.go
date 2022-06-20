package screen

import (
	"context"
	"log"
	"os"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	defColorFG tcell.Color // Default color foreground
	defColorBG tcell.Color // Default color background
	boxColorFG tcell.Color // Box color foreground
	boxColorBG tcell.Color // Box color background
	screen     tcell.Screen
	version    configs.LimeVersion
	cursorPos  Cursor // Cursor position
	getContext func() context.Context
	setContext func(context.Context)
	cancelCtx  context.CancelFunc
}

func NewScreen(version configs.LimeVersion, setCtx func(context.Context), getCtx func() context.Context, cancelCtx context.CancelFunc) *Screen {
	defStyle := utils.CreateStyle(version.DefaultBackgroundColor, version.DefaultForegroundColor)

	// Initialize screen
	tScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := tScreen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	tScreen.SetStyle(defStyle)
	tScreen.EnableMouse()
	tScreen.EnablePaste()
	tScreen.Clear()

	return &Screen{
		defColorFG: version.DefaultForegroundColor,
		defColorBG: version.DefaultBackgroundColor,
		boxColorFG: version.BoxForegroundColor,
		boxColorBG: version.BoxBackgroundColor,
		screen:     tScreen,
		version:    version,
		getContext: getCtx,
		setContext: setCtx,
		cancelCtx:  cancelCtx,
	}
}

func (s *Screen) Quit() {
	s.screen.Fini()
	os.Exit(0)
}

func (s *Screen) ShowBox() {
	s.GetScreen().Clear()
	// Draw background box
	sw, sh := s.screen.Size()     // screen width and height
	bw, bh := (sw/2)-30, (sh/2)-6 // box width and height
	s.DrawBox(bw, bh, sw-bw, sh-bh-3, "", true, true, true, tcell.StyleDefault)

	// Draw inner contents - logo, info, version, author...
	var (
		logoTop  = "█░░ █ █▀▄▀█ █▀▀"
		logoDown = "█▄▄ █ █░▀░█ ██▄"
		info1    = "Ctrl+W = start | Ctrl+Q = quit | Ctrl+S = save"
		info2    = "Ctrl+Z = undo | Ctrl+Y = redo"
	)

	ltx, lty := (sw/2)-8, (sh/2)-4              // logoTop x and y
	ldx, ldy := (sw/2)-8, (sh/2)-3              // logoDown x and y
	i1x, i1y := (sw/2)-len(info1)/2+1, (sh/2)-1 // info1 x and y
	i2x, i2y := (sw/2)-len(info2)/2+1, (sh/2)-0 // info2 x and y

	s.DrawText(ltx, lty, ltx+len(logoTop), lty+1, logoTop, s.GetBoxStyle())
	s.DrawText(ldx, ldy, ldx+len(logoDown), ldy+1, logoDown, s.GetBoxStyle())
	s.DrawText(i1x, i1y, i1x+len(info1), i1y+1, info1, s.GetBoxStyle())
	s.DrawText(i2x, i2y, i2x+len(info2), i2y+1, info2, s.GetBoxStyle())
}

func (s *Screen) GetScreen() tcell.Screen {
	return s.screen
}

func (s *Screen) GetDefStyle() tcell.Style {
	return utils.CreateStyle(s.defColorBG, s.defColorFG)
}

func (s *Screen) GetBoxStyle() tcell.Style {
	return utils.CreateStyle(s.boxColorBG, s.boxColorFG)
}

func (s *Screen) CancelContext() {
	s.cancelCtx()
}
