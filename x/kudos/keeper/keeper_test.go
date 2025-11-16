package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/keeper"
	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

// setupKeeper creates a keeper for testing
func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return k, ctx
}

func TestGetSetKudosBalance(t *testing.T) {
	k, ctx := setupKeeper(t)

	address := "cosmos1test"

	// Initially should be 0
	balance := k.GetKudosBalance(ctx, address)
	require.Equal(t, uint64(0), balance)

	// Set balance
	k.SetKudosBalance(ctx, address, 100)
	balance = k.GetKudosBalance(ctx, address)
	require.Equal(t, uint64(100), balance)
}

func TestAddKudos(t *testing.T) {
	k, ctx := setupKeeper(t)

	address := "cosmos1test"

	// Add kudos
	k.AddKudos(ctx, address, 50)
	balance := k.GetKudosBalance(ctx, address)
	require.Equal(t, uint64(50), balance)

	// Add more kudos
	k.AddKudos(ctx, address, 30)
	balance = k.GetKudosBalance(ctx, address)
	require.Equal(t, uint64(80), balance)
}

func TestSendKudos(t *testing.T) {
	k, ctx := setupKeeper(t)

	fromAddr := "cosmos1from"
	toAddr := "cosmos1to"

	// Send kudos
	err := k.SendKudos(ctx, fromAddr, toAddr, 100, "Great work!")
	require.NoError(t, err)

	// Check recipient balance
	balance := k.GetKudosBalance(ctx, toAddr)
	require.Equal(t, uint64(100), balance)

	// Check history counter
	counter := k.GetHistoryCounter(ctx)
	require.Equal(t, uint64(1), counter)
}

func TestDailyQuotaEnforcement(t *testing.T) {
	k, ctx := setupKeeper(t)

	fromAddr := "cosmos1from"
	toAddr := "cosmos1to"

	err := k.SendKudos(ctx, fromAddr, toAddr, types.DefaultDailyLimit, "using up quota")
	require.NoError(t, err)

	err = k.SendKudos(ctx, fromAddr, "cosmos1overflow", 1, "should exceed")
	require.ErrorIs(t, err, types.ErrDailyLimitExceeded)

	quota := k.GetDailyQuota(ctx, fromAddr)
	require.Equal(t, types.DefaultDailyLimit, quota.Used)
	require.Equal(t, uint64(0), quota.Remaining)
}

func TestDailyQuotaResetAfterWindow(t *testing.T) {
	k, ctx := setupKeeper(t)

	fromAddr := "cosmos1from"
	toAddr := "cosmos1to"

	require.NoError(t, k.SendKudos(ctx, fromAddr, toAddr, 10, "first"))

	ctx = ctx.WithBlockHeader(cmtproto.Header{Time: ctx.BlockTime().Add(time.Duration(types.DailyQuotaWindowSeconds+3600) * time.Second)})

	require.NoError(t, k.SendKudos(ctx, fromAddr, toAddr, types.DefaultDailyLimit, "after reset"))

	quota := k.GetDailyQuota(ctx, fromAddr)
	require.Equal(t, types.DefaultDailyLimit, quota.Used)
	require.Equal(t, uint64(0), quota.Remaining)
}

func TestSendKudosToSelf(t *testing.T) {
	k, ctx := setupKeeper(t)

	address := "cosmos1test"

	// Try to send kudos to self
	err := k.SendKudos(ctx, address, address, 100, "Self kudos")
	require.Error(t, err)
	require.Equal(t, types.ErrSameAddress, err)
}

func TestSendKudosZeroAmount(t *testing.T) {
	k, ctx := setupKeeper(t)

	fromAddr := "cosmos1from"
	toAddr := "cosmos1to"

	// Try to send zero kudos
	err := k.SendKudos(ctx, fromAddr, toAddr, 0, "Zero kudos")
	require.Error(t, err)
	require.Equal(t, types.ErrInvalidAmount, err)
}

func TestGetLeaderboard(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Set up multiple addresses with different balances
	k.SetKudosBalance(ctx, "cosmos1addr1", 100)
	k.SetKudosBalance(ctx, "cosmos1addr2", 200)
	k.SetKudosBalance(ctx, "cosmos1addr3", 50)

	// Get leaderboard
	leaderboard := k.GetLeaderboard(ctx, 10)
	require.Len(t, leaderboard, 3)

	// Check order (should be sorted by balance descending)
	require.Equal(t, "cosmos1addr2", leaderboard[0].Address)
	require.Equal(t, uint64(200), leaderboard[0].Balance)
	require.Equal(t, "cosmos1addr1", leaderboard[1].Address)
	require.Equal(t, uint64(100), leaderboard[1].Balance)
	require.Equal(t, "cosmos1addr3", leaderboard[2].Address)
	require.Equal(t, uint64(50), leaderboard[2].Balance)
}

func TestGetLeaderboardWithLimit(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Set up multiple addresses
	k.SetKudosBalance(ctx, "cosmos1addr1", 100)
	k.SetKudosBalance(ctx, "cosmos1addr2", 200)
	k.SetKudosBalance(ctx, "cosmos1addr3", 50)

	// Get leaderboard with limit
	leaderboard := k.GetLeaderboard(ctx, 2)
	require.Len(t, leaderboard, 2)
	require.Equal(t, "cosmos1addr2", leaderboard[0].Address)
	require.Equal(t, "cosmos1addr1", leaderboard[1].Address)
}
