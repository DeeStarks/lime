package highlighters

import (
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	INT_HIGHLIGHTER     = utils.CreateStyle(tcell.ColorReset, tcell.ColorSkyblue)
	STRING_HIGHLIGHTER  = utils.CreateStyle(tcell.ColorReset, tcell.ColorLightYellow)
	COMMENT_HIGHLIGHTER = utils.CreateStyle(tcell.ColorReset, tcell.ColorDimGray)
)

type Highlighter struct {
	ext              string // The file extension
	scheme           map[string]tcell.Style
	commentIndicator string
}

func NewHighlighter(ext string) *Highlighter {
	var (
		scheme           map[string]tcell.Style
		commentIndicator string
	)

	switch ext {
	case ".go":
		scheme = golang_scheme
		commentIndicator = golang_comment_indicator
	case ".py":
		scheme = python_scheme
		commentIndicator = python_comment_indicator
	}

	return &Highlighter{
		ext:              ext,
		scheme:           scheme,
		commentIndicator: commentIndicator,
	}
}

func (h *Highlighter) GetStyle(word string) tcell.Style {
	if style, ok := h.scheme[word]; ok {
		return style
	}
	return utils.CreateStyle(tcell.ColorReset, tcell.ColorWhite)
}

func (h *Highlighter) GetCommentIndicator() string {
	return h.commentIndicator
}
