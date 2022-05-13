package screen

import "github.com/gdamore/tcell/v2"

func (s *Screen) DrawText(x1, y1, x2, y2 int, text string, style tcell.Style) {
	tScreen := s.GetScreen()

	row := y1
	col := x1
	for _, r := range text {
		if style == tcell.StyleDefault {
			tScreen.SetContent(col, row, r, nil, s.GetBoxStyle())
		} else {
			tScreen.SetContent(col, row, r, nil, style)
		}
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (s *Screen) DrawBox(x1, y1, x2, y2 int, text string, fill, borders, corners bool, style tcell.Style) {
	tScreen := s.GetScreen() //

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	if fill {
		for row := y1; row <= y2; row++ {
			for col := x1; col <= x2; col++ {
				if style == tcell.StyleDefault {
					tScreen.SetContent(col, row, ' ', nil, s.GetBoxStyle())
				} else {
					tScreen.SetContent(col, row, ' ', nil, style)
				}
			}
		}
	}

	// Draw borders
	if borders {
		for col := x1; col <= x2; col++ {
			tScreen.SetContent(col, y1, tcell.RuneHLine, nil, s.GetBoxStyle())
			tScreen.SetContent(col, y2, tcell.RuneHLine, nil, s.GetBoxStyle())
		}
		for row := y1 + 1; row < y2; row++ {
			tScreen.SetContent(x1, row, tcell.RuneVLine, nil, s.GetBoxStyle())
			tScreen.SetContent(x2, row, tcell.RuneVLine, nil, s.GetBoxStyle())
		}
	}

	// Only draw corners if necessary
	if corners {
		if y1 != y2 && x1 != x2 {
			tScreen.SetContent(x1, y1, tcell.RuneULCorner, nil, s.GetBoxStyle())
			tScreen.SetContent(x2, y1, tcell.RuneURCorner, nil, s.GetBoxStyle())
			tScreen.SetContent(x1, y2, tcell.RuneLLCorner, nil, s.GetBoxStyle())
			tScreen.SetContent(x2, y2, tcell.RuneLRCorner, nil, s.GetBoxStyle())
		}
	}

	if text != "" {
		s.DrawText(x1+1, y1+1, x2-1, y2-1, text, style)
	}
}
