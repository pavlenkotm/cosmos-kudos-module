package keeper

import (
	"encoding/binary"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger
}

// NewKeeper creates a new kudos Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetKudosBalance returns the kudos balance for an address
func (k Keeper) GetKudosBalance(ctx sdk.Context, address string) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	key := types.KudosBalanceKey(address)

	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

// SetKudosBalance sets the kudos balance for an address
func (k Keeper) SetKudosBalance(ctx sdk.Context, address string, balance uint64) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.KudosBalanceKey(address)

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, balance)

	if err := store.Set(key, bz); err != nil {
		panic(err)
	}
}

// AddKudos adds kudos to an address balance
func (k Keeper) AddKudos(ctx sdk.Context, address string, amount uint64) {
	currentBalance := k.GetKudosBalance(ctx, address)
	newBalance := currentBalance + amount
	k.SetKudosBalance(ctx, address, newBalance)
}

// GetHistoryCounter returns the current history counter
func (k Keeper) GetHistoryCounter(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(types.HistoryCounterKey)
	if err != nil || bz == nil {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

// SetHistoryCounter sets the history counter
func (k Keeper) SetHistoryCounter(ctx sdk.Context, counter uint64) {
	store := k.storeService.OpenKVStore(ctx)

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, counter)

	if err := store.Set(types.HistoryCounterKey, bz); err != nil {
		panic(err)
	}
}

// AddKudosHistory adds a kudos transaction to history
func (k Keeper) AddKudosHistory(ctx sdk.Context, fromAddress, toAddress string, amount uint64, comment string) {
	counter := k.GetHistoryCounter(ctx)
	counter++

	history := types.KudosHistory{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
		Comment:     comment,
		Timestamp:   time.Now().Unix(),
	}

	store := k.storeService.OpenKVStore(ctx)
	key := types.KudosHistoryKey(counter)
	bz := k.cdc.MustMarshal(&history)

	if err := store.Set(key, bz); err != nil {
		panic(err)
	}

	k.SetHistoryCounter(ctx, counter)
}

// GetAllKudosBalances returns all kudos balances
func (k Keeper) GetAllKudosBalances(ctx sdk.Context) map[string]uint64 {
	store := k.storeService.OpenKVStore(ctx)
	balances := make(map[string]uint64)

	iterator, err := store.Iterator(types.KudosBalancePrefix, sdk.PrefixEndBytes(types.KudosBalancePrefix))
	if err != nil {
		return balances
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		address := string(key[len(types.KudosBalancePrefix):])
		balance := binary.BigEndian.Uint64(iterator.Value())
		balances[address] = balance
	}

	return balances
}

// GetLeaderboard returns the top N kudos receivers
func (k Keeper) GetLeaderboard(ctx sdk.Context, limit uint32) []types.LeaderboardEntry {
	balances := k.GetAllKudosBalances(ctx)

	// Convert map to slice for sorting
	entries := make([]types.LeaderboardEntry, 0, len(balances))
	for address, balance := range balances {
		entries = append(entries, types.LeaderboardEntry{
			Address: address,
			Balance: balance,
		})
	}

	// Sort by balance descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Balance > entries[j].Balance
	})

	// Apply limit
	if limit > 0 && uint32(len(entries)) > limit {
		entries = entries[:limit]
	}

	return entries
}

// SendKudos sends kudos from one address to another
func (k Keeper) SendKudos(ctx sdk.Context, fromAddress, toAddress string, amount uint64, comment string) error {
	// Validate addresses are different
	if fromAddress == toAddress {
		return types.ErrSameAddress
	}

	// Validate amount
	if amount == 0 {
		return types.ErrInvalidAmount
	}

	// Add kudos to recipient
	k.AddKudos(ctx, toAddress, amount)

	// Add to history
	k.AddKudosHistory(ctx, fromAddress, toAddress, amount, comment)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.ModuleName,
			sdk.NewAttribute("action", "send_kudos"),
			sdk.NewAttribute("from", fromAddress),
			sdk.NewAttribute("to", toAddress),
			sdk.NewAttribute("amount", fmt.Sprintf("%d", amount)),
		),
	)

	return nil
}
