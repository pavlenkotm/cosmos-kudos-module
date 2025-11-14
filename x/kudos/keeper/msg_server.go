package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the kudos MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// SendKudos implements the SendKudos message handler
func (k msgServer) SendKudos(goCtx context.Context, msg *types.MsgSendKudos) (*types.MsgSendKudosResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Send kudos
	if err := k.Keeper.SendKudos(ctx, msg.FromAddress, msg.ToAddress, msg.Amount, msg.Comment); err != nil {
		return nil, err
	}

	return &types.MsgSendKudosResponse{}, nil
}
