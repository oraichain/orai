package types

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		AIOracles: []AIOracle{},
		Params:    DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// ValidateGenesis validates the aioracle genesis parameters
func (gs *GenesisState) Validate() error {
	// for _, record := range gs.AIOracles {
	// 	if record.RequestID == "" {
	// 		return fmt.Errorf("invalid aiOracles: Value: %s. Error: Missing RequestID", record.RequestID)
	// 	}
	// 	if record.Creator == nil {
	// 		return fmt.Errorf("invalid aiOracles: Owner: %s. Error: Missing Creator", record.Creator)
	// 	}
	// 	// if record.Fees.Empty() {
	// 	// 	return fmt.Errorf("invalid aiOracles: Owner: %s. Error: Missing Fees", record.Fees)
	// 	// }
	// }
	return nil
}
