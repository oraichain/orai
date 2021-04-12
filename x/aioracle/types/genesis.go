package types

import fmt "fmt"

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 2

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		AIOracles:       []AIOracle{},
		Reports:         []Report{},
		AIOracleResults: []AIOracleResult{},
		Params:          DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// ValidateGenesis validates the aioracle genesis parameters
func (gs *GenesisState) Validate() error {
	for _, record := range gs.AIOracles {
		if record.RequestID == "" {
			return fmt.Errorf("invalid aioracles: Value: %s. Error: Missing RequestID", record.RequestID)
		}
		if record.Creator == nil {
			return fmt.Errorf("invalid aioracles: Owner: %s. Error: Missing Creator", record.Creator)
		}
		if record.Contract.Empty() {
			return fmt.Errorf("invalid aioracles: Owner: %s. Error: Missing Oracle Script Contract")
		}
		// if record.Fees.Empty() {
		// 	return fmt.Errorf("invalid aioracles: Owner: %s. Error: Missing Fees", record.Fees)
		// }
	}
	return nil
}
