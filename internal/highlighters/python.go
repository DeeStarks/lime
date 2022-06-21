package highlighters

import (
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	python_comment_indicator = "#"
	// Highlighters is a map of highlighters
	python_scheme = map[string]tcell.Style{
		"class":    utils.CreateStyle(tcell.ColorBlack, tcell.ColorBlue),
		"def":      utils.CreateStyle(tcell.ColorBlack, tcell.ColorBlue),
		"global":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorBlue),
		"nonlocal": utils.CreateStyle(tcell.ColorBlack, tcell.ColorBlue),
		"assert":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorBlue),

		"not": utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		"and": utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		"or":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		"=":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		"==":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		"!=":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		"<":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		">":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		"<=":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),
		">=":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorViolet),

		"import":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"from":     utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"as":       utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"if":       utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"elif":     utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"else":     utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"while":    utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"for":      utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"in":       utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"try":      utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"except":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"finally":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"with":     utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"pass":     utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"continue": utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"break":    utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"return":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"lambda":   utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),
		"->":       utils.CreateStyle(tcell.ColorBlack, tcell.ColorPurple),

		"True":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorDeepSkyBlue),
		"False": utils.CreateStyle(tcell.ColorBlack, tcell.ColorDeepSkyBlue),
		"None":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorDeepSkyBlue),

		"self":     utils.CreateStyle(tcell.ColorBlack, tcell.ColorSlateBlue),
		"print":    utils.CreateStyle(tcell.ColorBlack, tcell.ColorSlateBlue),
		"__init__": utils.CreateStyle(tcell.ColorBlack, tcell.ColorSlateBlue),
		"__str__":  utils.CreateStyle(tcell.ColorBlack, tcell.ColorSlateBlue),
	}
)
