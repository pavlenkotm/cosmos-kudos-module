package types

// GenesisState defines the kudos module's genesis state.
type GenesisState struct {
	// Add genesis fields if needed in the future
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState() *GenesisState {
	return &GenesisState{}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState()
}
