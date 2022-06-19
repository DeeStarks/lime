package highlighters

import (
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	INT_HIGHLIGHTER     = utils.CreateStyle(tcell.ColorReset, tcell.ColorSkyblue)
	STRING_HIGHLIGHTER  = utils.CreateStyle(tcell.ColorReset, tcell.ColorLightGoldenrodYellow)
	COMMENT_HIGHLIGHTER = utils.CreateStyle(tcell.ColorReset, tcell.ColorDimGray)
)

type Highlighter struct {
	ext    string // The file extension
	scheme map[string]tcell.Style
}

func NewHighlighter(ext string) *Highlighter {
	var scheme map[string]tcell.Style

	switch ext {
	case ".go":
		scheme = golang_scheme
	case ".py":
		scheme = python_scheme
	}

	return &Highlighter{
		ext:    ext,
		scheme: scheme,
	}
}

func (h *Highlighter) GetStyle(word string) tcell.Style {
	if style, ok := h.scheme[word]; ok {
		return style
	}
	return utils.CreateStyle(tcell.ColorReset, tcell.ColorWhite)
}
