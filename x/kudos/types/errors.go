package types

import (
	"cosmossdk.io/errors"
)

// Kudos module sentinel errors
var (
	ErrInvalidAddress     = errors.Register(ModuleName, 1, "invalid address")
	ErrSameAddress        = errors.Register(ModuleName, 2, "cannot send kudos to yourself")
	ErrInvalidAmount      = errors.Register(ModuleName, 3, "amount must be greater than 0")
	ErrCommentTooLong     = errors.Register(ModuleName, 4, "comment exceeds 140 characters")
	ErrInvalidLeaderboard = errors.Register(ModuleName, 5, "invalid leaderboard parameters")
	ErrDailyLimitExceeded = errors.Register(ModuleName, 6, "daily kudos limit exceeded")
)
