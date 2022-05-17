package screen

import (
	"context"

	"github.com/DeeStarks/lime/internal/constants"
	"github.com/DeeStarks/lime/internal/utils"
)

func (s *Screen) ScrollUp() {
	// Decrement the startline
	sl := s.getContext().Value(constants.StartLineCtxKey).(int)
	if sl > 0 {
		ctx := context.WithValue(s.getContext(), constants.StartLineCtxKey, sl-1)
		utils.LogMessage("%s", sl+1)
		s.setContext(ctx)
	}
}

func (s *Screen) ScrollDown() {
	// Increment the startline
	sl := s.getContext().Value(constants.StartLineCtxKey).(int)
	ctx := context.WithValue(s.getContext(), constants.StartLineCtxKey, sl+1)
	utils.LogMessage("%s", sl+1)
	s.setContext(ctx)
}
