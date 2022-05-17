package screen

import (
	"context"

	"github.com/DeeStarks/lime/internal/constants"
)

func (s *Screen) ScrollUp() {
	// Decrement the startline
	sl := s.getContext().Value(constants.StartLineCtxKey).(int)
	if sl > 0 {
		ctx := context.WithValue(s.getContext(), constants.StartLineCtxKey, sl-1)
		s.setContext(ctx)
	}
}

func (s *Screen) ScrollDown() {
	// Increment the startline
	sl := s.getContext().Value(constants.StartLineCtxKey).(int)
	_, sh := s.GetScreen().Size()
	sh = sh - constants.EditorPaddingTop - constants.EditorPaddingBottom
	tl := s.getContext().Value(constants.TotalLinesCtxKey).(int)
	if sl+sh < tl {
		ctx := context.WithValue(s.getContext(), constants.StartLineCtxKey, sl+1)
		s.setContext(ctx)
	}
}
