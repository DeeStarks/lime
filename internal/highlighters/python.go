package highlighters

import (
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	python_comment_indicator = "#"
	// Highlighters is a map of highlighters
	python_scheme = map[string]tcell.Style{
		"class":    utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"def":      utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"global":   utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"nonlocal": utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"assert":   utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),

		"not": utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		"and": utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		"or":  utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		"=":   utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		"==":  utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		"!=":  utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		"<":   utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		">":   utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		"<=":  utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),
		">=":  utils.CreateStyle(tcell.ColorReset, tcell.ColorViolet),

		"import":   utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"from":     utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"as":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"if":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"elif":     utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"else":     utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"while":    utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"for":      utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"in":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"try":      utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"except":   utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"finally":  utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"with":     utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"pass":     utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"continue": utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"break":    utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"return":   utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"lambda":   utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"->":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),

		"True":  utils.CreateStyle(tcell.ColorReset, tcell.ColorDeepSkyBlue),
		"False": utils.CreateStyle(tcell.ColorReset, tcell.ColorDeepSkyBlue),
		"None":  utils.CreateStyle(tcell.ColorReset, tcell.ColorDeepSkyBlue),

		"self":     utils.CreateStyle(tcell.ColorReset, tcell.ColorSlateBlue),
		"print":    utils.CreateStyle(tcell.ColorReset, tcell.ColorSlateBlue),
		"__init__": utils.CreateStyle(tcell.ColorReset, tcell.ColorSlateBlue),
		"__str__":  utils.CreateStyle(tcell.ColorReset, tcell.ColorSlateBlue),
	}
)
