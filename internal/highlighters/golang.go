package highlighters

import (
	"github.com/DeeStarks/lime/internal/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	golang_comment_indicator = "//"
	// Highlighters is a map of highlighters
	golang_scheme = map[string]tcell.Style{
		"func":        utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"var":         utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"const":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"type":        utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"struct":      utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"interface":   utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"package":     utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"import":      utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"return":      utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"break":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"continue":    utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"goto":        utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"fallthrough": utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"if":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"else":        utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"for":         utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"range":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"switch":      utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"case":        utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"default":     utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"select":      utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"defer":       utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"go":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		":=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"=":           utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"!=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"==":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"<":           utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		">":           utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"<=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		">=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"+=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"-=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"*=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"/=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"%=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"&=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"|=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"^=":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"<<=":         utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		">>=":         utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"&&":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),
		"||":          utils.CreateStyle(tcell.ColorReset, tcell.ColorPurple),

		"int":        utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"string":     utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"bool":       utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"float64":    utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"float32":    utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"uint":       utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"uint8":      utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"uint16":     utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"uint32":     utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"uint64":     utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"uintptr":    utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"byte":       utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"rune":       utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"complex64":  utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),
		"complex128": utils.CreateStyle(tcell.ColorReset, tcell.ColorBlue),

		"true":  utils.CreateStyle(tcell.ColorReset, tcell.ColorSkyblue),
		"false": utils.CreateStyle(tcell.ColorReset, tcell.ColorSkyblue),
		"nil":   utils.CreateStyle(tcell.ColorReset, tcell.ColorSkyblue),
		"iota":  utils.CreateStyle(tcell.ColorReset, tcell.ColorSkyblue),
	}
)
