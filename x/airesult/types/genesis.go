package types

import fmt "fmt"

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 2

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		AIRequestResults: []AIRequestResult{},
		Params:           DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// Validate validates the provider genesis parameters
func (gs GenesisState) Validate() error {
	for _, record := range gs.AIRequestResults {
		if record.RequestID == "" {
			return fmt.Errorf("invalid AIRequestResults: Value: %s. Error: Missing RequestID", record.RequestID)
		}
		if len(record.Status) == 0 {
			return fmt.Errorf("invalid AIRequestResults: Status: %s. Error: Missing Status", record.Status)
		}
	}
	return nil
}
