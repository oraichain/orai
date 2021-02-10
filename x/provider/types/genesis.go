package types

import fmt "fmt"

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		OracleScripts: []OracleScript{},
		AIDataSources: []AIDataSource{},
		TestCases:     []TestCase{},
		Params:        DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs *GenesisState) Validate() error {
	for _, record := range gs.OracleScripts {
		if record.Name == "" {
			return fmt.Errorf("invalid OracleScripts: Value: %s. Error: Missing Name", record.Name)
		}
		if record.Owner == nil {
			return fmt.Errorf("invalid OracleScripts: Owner: %s. Error: Missing Owner", record.Owner)
		}
	}

	for _, record := range gs.AIDataSources {
		if record.Name == "" {
			return fmt.Errorf("invalid AIDataSources: Value: %s. Error: Missing Name", record.Name)
		}
		if record.Owner == nil {
			return fmt.Errorf("invalid AIDataSources: Owner: %s. Error: Missing Owner", record.Owner)
		}
	}

	for _, record := range gs.TestCases {
		if len(record.Owner) == 0 {
			return fmt.Errorf("invalid request: Owner: %s. Error: Missing requestID", string(record.Owner[:]))
		}
		if len(record.Name) == 0 {
			return fmt.Errorf("invalid request: Name: %s. Error: Missing name", record.Name)
		}
	}
	return nil
}
