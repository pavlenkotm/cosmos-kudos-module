package types

const (
	// ModuleName is the name of the kudos module
	ModuleName = "kudos"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the kudos module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the kudos module
	QuerierRoute = ModuleName
)

var (
	// KudosBalancePrefix is the prefix for kudos balance keys
	KudosBalancePrefix = []byte{0x01}

	// KudosHistoryPrefix is the prefix for kudos history keys
	KudosHistoryPrefix = []byte{0x02}

	// HistoryCounterKey is the key for the global history counter
	HistoryCounterKey = []byte{0x03}
)

// KudosBalanceKey returns the key for a kudos balance
func KudosBalanceKey(address string) []byte {
	return append(KudosBalancePrefix, []byte(address)...)
}

// KudosHistoryKey returns the key for a kudos history entry
func KudosHistoryKey(id uint64) []byte {
	bz := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bz[i] = byte(id >> (8 * (7 - i)))
	}
	return append(KudosHistoryPrefix, bz...)
}
