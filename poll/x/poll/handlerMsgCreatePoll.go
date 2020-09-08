package poll

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skhiearth/poll/x/poll/types"
	"github.com/skhiearth/poll/x/poll/keeper"
)

func handleMsgCreatePoll(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreatePoll) (*sdk.Result, error) {
	var poll = types.Poll{
		Creator: msg.Creator,
		ID:      msg.ID,
    Title: msg.Title,
    Options: msg.Options,
	}
	k.CreatePoll(ctx, poll)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
