package editor

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/constants"
)

func (e *Editor) Read() {
	// Clear the editor
	e.Clear()

	screen, _ := drawEditBox(e)
	sw, _ := e.screen.GetScreen().Size()
	
	buf := e.ReadBufferByte()

	tab := make([]byte, configs.TabSize)
	for i := 0; i < configs.TabSize; i++ {
		tab[i] = ' '
	}

	var (
		word string
		wordCount int
		xAxis int = constants.EditorPaddingLeft + 2
		yAxis int = constants.EditorPaddingTop + 1
		lineLength int = 0
		lineCounter []int
		numCount int = 1
		numbering = []LineNumbering{
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
				word, e.defStyle)
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
		} else if c == ' ' {
			lineLength++
			if xAxis + wordCount > sw {
				xAxis = constants.EditorPaddingLeft + 2
				yAxis++
				lineCounter = append(lineCounter, lineLength-wordCount)
				lineLength = 0
			}
			e.screen.DrawText(
				xAxis, yAxis, 
				xAxis+wordCount, yAxis,
				word, e.defStyle)
			xAxis += wordCount+1
			word = ""
			wordCount = 0
		} else if i == len(buf)-1 {
			word += string(c)
			wordCount++
			lineLength++
			if xAxis + wordCount > sw {
				xAxis = constants.EditorPaddingLeft + 2
				yAxis++
			}
			e.screen.DrawText(
				xAxis, yAxis,
				xAxis+wordCount, yAxis,
				word, e.defStyle)
		} else {
			word += string(c)
			wordCount++
			lineLength++
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
