package types

import fmt "fmt"

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 2

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		AIRequests: []AIRequest{},
		Params:     DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// ValidateGenesis validates the airequest genesis parameters
func (gs *GenesisState) Validate() error {
	for _, record := range gs.AIRequests {
		if record.RequestID == "" {
			return fmt.Errorf("invalid AIRequests: Value: %s. Error: Missing RequestID", record.RequestID)
		}
		if record.Creator == nil {
			return fmt.Errorf("invalid AIRequests: Owner: %s. Error: Missing Creator", record.Creator)
		}
		if record.OracleScriptName == "" {
			return fmt.Errorf("invalid AIRequests: Owner: %s. Error: Missing Oracle Script Name", record.OracleScriptName)
		}
		// if record.Fees.Empty() {
		// 	return fmt.Errorf("invalid AIRequests: Owner: %s. Error: Missing Fees", record.Fees)
		// }
	}
	return nil
}
