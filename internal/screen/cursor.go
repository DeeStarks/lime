package screen

import (
	"context"

	"github.com/DeeStarks/lime/internal/constants"
)

type Cursor struct {
	x, y int
}

func (s *Screen) SetCursor(x, y int) {
	s.cursorPos.x = x
	s.cursorPos.y = y
	s.ShowCursor()

	// Update buffer index. This is used to determine the index to write to
	index := y*(x-constants.EditorPaddingLeft-2)-1 // -1 to begin at 0
	ctx := context.WithValue(s.getContext(), constants.BufferIndexCtxKey, index)
	s.setContext(ctx)
}

func (s *Screen) GetCursorPosition() (x int, y int) {
	x, y = s.cursorPos.x, s.cursorPos.y
	return
}

func (s *Screen) ShowCursor() {
	s.GetScreen().ShowCursor(s.cursorPos.x, s.cursorPos.y)
}

func (s *Screen) MoveCursor(code constants.KeyCode) {
	switch code {
	case constants.KeyArrowLeft:
		x, y := s.GetCursorPosition()
		pl := constants.EditorPaddingLeft + 2 // +2 for editor's extrapadding
		if x-1 <= pl {
			s.SetCursor(pl, y)
		} else {
			s.SetCursor(x-1, y)
		}
	case constants.KeyArrowRight:
		x, y := s.GetCursorPosition()
		sw, _ := s.GetScreen().Size()
		pr := constants.EditorPaddingRight + 1
		if x+1 >= sw-pr {
			s.SetCursor(sw-pr, y)
		} else {
			s.SetCursor(x+1, y)
		}
	case constants.KeyArrowUp:
		x, y := s.GetCursorPosition()
		pt := constants.EditorPaddingTop + 1
		if y-1 <= pt {
			s.ScrollUp()
			s.SetCursor(x, pt)
		} else {
			s.SetCursor(x, y-1)
		}
	case constants.KeyArrowDown:
		x, y := s.GetCursorPosition()
		_, sh := s.GetScreen().Size()
		pb := constants.EditorPaddingBottom + 1
		if y+1 >= sh-pb {
			s.ScrollDown()
			s.SetCursor(x, sh-pb)
		} else {
			s.SetCursor(x, y+1)
		}
	case constants.KeyEnter:
		x, y := s.GetCursorPosition()
		_, sh := s.GetScreen().Size()
		pb, pl := constants.EditorPaddingBottom+1, constants.EditorPaddingLeft+2
		if y == sh-pb {
			s.ScrollDown()
			s.SetCursor(pl, sh-pb)
		} else {
			s.SetCursor(x, y+1)
		}
		s.SetCursor(x, y+1)
	}
}
