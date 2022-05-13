package utils

import "github.com/gdamore/tcell/v2"

func CreateStyle(bg, fg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}
