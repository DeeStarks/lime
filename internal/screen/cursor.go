package screen

type Cursor struct {
	x, y int
}

func (s *Screen) SetCursor(x, y int) {
	s.cursorPos.x = x
	s.cursorPos.y = y
	s.ShowCursor()
}

func (s *Screen) GetCursorPosition() (x int, y int) {
	x, y = s.cursorPos.x, s.cursorPos.y
	return
}

func (s *Screen) ShowCursor() {
	s.GetScreen().ShowCursor(s.cursorPos.x, s.cursorPos.y)
}
