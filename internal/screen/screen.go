package screen

import (
	"fmt"
	"log"
	"os"

	"github.com/DeeStarks/lime/configs"
	"github.com/gdamore/tcell/v2"
)

func createStyle(bg, fg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}

type Screen struct {
	defColorFG tcell.Color // Default color foreground
	defColorBG tcell.Color // Default color background
	boxColorFG tcell.Color // Box color foreground
	boxColorBG tcell.Color // Box color background
	screen     tcell.Screen
	version    configs.LimeVersion
}

func NewScreen(version configs.LimeVersion) *Screen {
	defStyle := createStyle(version.DefaultBackgroundColor, version.DefaultForegroundColor)

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
	}
}

func (s *Screen) Quit() {
	s.screen.Fini()
	os.Exit(0)
}

func (s *Screen) ShowBox() {
	// Draw background box
	sw, sh := s.screen.Size()     // screen width and height
	bw, bh := (sw/2)-30, (sh/2)-5 // box width and height
	s.DrawBox(bw, bh, sw-bw, sh-bh-3, "", true, true, true)

	// Draw inner contents - logo, info, version, author...
	// Logo
	lw, lh := (sw/2)-9, (sh/2)-3 // Logo size
	logo := s.version.Logo
	// The logo is displayed in reverse, so we need to reverse it
	left, right := 0, len(logo)-1
	for left < right {
		logo[left], logo[right] = logo[right], logo[left]
		left++
		right--
	}

	for _, v := range s.version.Logo {
		fmt.Println(v)
		s.DrawBox(lw, lh, sw-lw, sh-lh, v, false, false, false)
		lh-- // Subtract to enter next line
	}
	// Info text
	iw, ih := (sw/2)-(len(s.version.InfoText)/2)-2, (sh/2)+1
	s.DrawBox(iw, ih, sw-iw, sh-ih, s.version.InfoText, false, false, false)

}

func (s *Screen) GetScreen() tcell.Screen {
	return s.screen
}

func (s *Screen) GetDefStyle() tcell.Style {
	return createStyle(s.defColorBG, s.defColorFG)
}

func (s *Screen) GetBoxStyle() tcell.Style {
	return createStyle(s.boxColorBG, s.boxColorFG)
}
