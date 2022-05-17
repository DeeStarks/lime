package constants

import "github.com/DeeStarks/lime/internal/global"

const (
	LineCounterCtxKey = global.ContextKey("line_counter")
	BufferIndexCtxKey = global.ContextKey("buffer_index")
	StartLineCtxKey   = global.ContextKey("start_line")
	TotalLinesCtxKey  = global.ContextKey("total_lines")
)
