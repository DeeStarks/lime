package editor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/highlighters"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) Read() {
	// Clear the editor
	e.Clear()

	screen, _ := drawEditBox(e)
	sw, _ := e.screen.GetScreen().Size()

	buf := e.ReadBufferByte()
	
	newlineCounts := 1
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\n' {
			if newlineCounts == e.startLine {
				buf = buf[i+1:]
				break
			}
			newlineCounts++
		}
	}

	tab := make([]byte, configs.TabSize)
	for i := 0; i < configs.TabSize; i++ {
		tab[i] = ' '
	}

	var (
		word        string
		wordCount   int
		wordStyle   tcell.Style
		openComment bool
		xAxis       int = constants.EditorPaddingLeft + 2
		yAxis       int = constants.EditorPaddingTop + 1
		lineLength  int = 0
		lineCounter []int
		numCount    int = e.startLine+1
		numbering       = []LineNumbering{
			{
				number: numCount,
				yAxis:  yAxis,
			},
		}
	)

	for i, c := range buf {
		if xAxis > e.width {
			xAxis = constants.EditorPaddingLeft + 2
			yAxis++
			lineCounter = append(lineCounter, lineLength)
			lineLength = 0
		}
		if c == '\t' {
			e.screen.DrawText(
				xAxis, yAxis,
				xAxis+configs.TabSize, yAxis,
				string(tab), e.defStyle)
			xAxis += configs.TabSize
			lineLength += configs.TabSize
		} else if c == '\n' {
			e.screen.DrawText(
				xAxis, yAxis,
				xAxis+wordCount, yAxis,
				word, wordStyle)
			xAxis = constants.EditorPaddingLeft + 2
			yAxis++
			word = ""
			wordCount = 0
			lineCounter = append(lineCounter, lineLength)
			lineLength = 0

			// Append numbering
			numCount++
			nmb := struct {
				number int
				yAxis  int
			}{
				number: numCount,
				yAxis:  yAxis,
			}
			numbering = append(numbering, nmb)

			openComment = false
			wordStyle = e.defStyle
		} else if c == ' ' {
			lineLength++
			if xAxis+wordCount > sw {
				xAxis = constants.EditorPaddingLeft + 2
				yAxis++
				lineCounter = append(lineCounter, lineLength-wordCount)
				lineLength = 0
			}
			e.screen.DrawText(
				xAxis, yAxis,
				xAxis+wordCount, yAxis,
				word, wordStyle)
			xAxis += wordCount + 1
			word = ""
			wordCount = 0
		} else if i == len(buf)-1 {
			word += string(c)
			wordCount++
			lineLength++
			if xAxis+wordCount > sw {
				xAxis = constants.EditorPaddingLeft + 2
				yAxis++
			}
			e.screen.DrawText(
				xAxis, yAxis,
				xAxis+wordCount, yAxis,
				word, wordStyle)
		} else {
			word += string(c)
			wordCount++
			lineLength++
		}

		// Highlight the word
		commentIndicator := e.highlighters.GetCommentIndicator()
		if len(word) > len(commentIndicator)-1 && word[0:len(commentIndicator)] == commentIndicator {
			openComment = true
			wordStyle = highlighters.COMMENT_HIGHLIGHTER
		}

		if !openComment {
			if len(word) > 0 && (word[0] == '"' && word[len(word)-1] == '"') {
				wordStyle = highlighters.STRING_HIGHLIGHTER
			} else if len(word) > 0 && (word[0] == '`' && word[len(word)-1] == '`') {
				wordStyle = highlighters.STRING_HIGHLIGHTER
			} else if _, ok := strconv.Atoi(word); ok == nil {
				wordStyle = highlighters.INT_HIGHLIGHTER
			} else if _, ok := strconv.ParseFloat(word, 64); ok == nil {
				wordStyle = highlighters.INT_HIGHLIGHTER
			} else {
				wordStyle = e.highlighters.GetStyle(word)
			}
		}
	}

	if e.ReadBufferString() != e.ReadInitialBufferString() {
		e.title = fmt.Sprintf("•• %s (modified) ••", e.file.Name())
	} else {
		e.title = fmt.Sprintf("•• %s ••", e.file.Name())
	}

	lineCounter = append(lineCounter, lineLength)
	e.lineCounter = lineCounter
	e.lineNumbering = numbering
	screen.Show()
	e.showBars()
}

func (e *Editor) Clear() {
	sw, sh := e.screen.GetScreen().Size()
	ew, eh := sw-constants.EditorPaddingLeft-constants.EditorPaddingRight, sh-constants.EditorPaddingTop-constants.EditorPaddingBottom

	// Create a new buffer filled with spaces and draw on the editor
	buf := make([]byte, ew)
	for i := 0; i < ew; i++ {
		buf[i] = ' '
	}

	// Draw the content
	for line := constants.EditorPaddingTop + 1; line <= eh; line++ {
		e.screen.DrawText(constants.EditorPaddingLeft, line, sw-constants.EditorPaddingRight, line+1, string(buf), e.defStyle)
	}
}

func (e *Editor) Write(char rune) {
	e.InsertToBuffer([]byte(string(char)), e.insIndex)
	e.insIndex++
	e.Read()
	e.UpdateCursorPosition()
}

func (e *Editor) BackSpace() {
	if e.insIndex > 0 {
		e.insIndex--
		e.RemoveFromBuffer(e.insIndex)
		e.UpdateCursorPosition()
		e.Read()
	}
}

func (e *Editor) Save() {
	bt := e.ReadBufferByte()
	e.initialBuffer = bytes.NewBuffer(bt)
	ioutil.WriteFile(e.file.Name(), e.ReadInitialBufferByte(), 0644)
	e.Read() // To update title bar
}
