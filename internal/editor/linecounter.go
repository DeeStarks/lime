package editor

type LineCounter []int // Each element represents a line, and the value is the number of characters

func NewLineCounter(v int) LineCounter {
	c := make(LineCounter, 1)
	c[0] = v
	return c
}
