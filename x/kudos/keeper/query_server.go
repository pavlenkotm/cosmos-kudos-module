package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

var _ types.QueryServer = Keeper{}

// KudosBalance implements the Query/KudosBalance gRPC method
func (k Keeper) KudosBalance(goCtx context.Context, req *types.QueryKudosBalanceRequest) (*types.QueryKudosBalanceResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidAddress
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get balance
	balance := k.GetKudosBalance(ctx, req.Address)

	return &types.QueryKudosBalanceResponse{
		Balance: balance,
	}, nil
}

// KudosLeaderboard implements the Query/KudosLeaderboard gRPC method
func (k Keeper) KudosLeaderboard(goCtx context.Context, req *types.QueryKudosLeaderboardRequest) (*types.QueryKudosLeaderboardResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidLeaderboard
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Default limit if not specified
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}

	// Get leaderboard
	entries := k.GetLeaderboard(ctx, limit)

	return &types.QueryKudosLeaderboardResponse{
		Entries: entries,
	}, nil
}
