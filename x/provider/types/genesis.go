package types

import (
	"fmt"
)

// GenesisState - all provider state that must be provided at genesis
type GenesisState struct {
	OracleScripts []OracleScript `json:"oracle_scripts"`
	AIDataSources []AIDataSource `json:"datasources"`
	TestCases     []TestCase     `json:"test_cases"`
	Params        Params         `json:"params"`
	// TODO: Fill out what is needed by the module for genesis
}

// // NewGenesisState creates a new GenesisState object
// func NewGenesisState(oScripts []OracleScript, aiDataSources []AIDataSource, aiRequests []AIRequest, testCases []TestCase, params Params) GenesisState {
// 	return GenesisState{
// 		// TODO: Fill out according to your genesis state
// 		OracleScripts: oScripts,
// 		AIDataSources: aiDataSources,
// 		TestCases:     testCases,
// 		Params:        params,
// 	}
// }

// NewGenesisState creates a new GenesisState object
func NewGenesisState(oScripts []OracleScript, aiDataSources []AIDataSource, testCases []TestCase, params Params) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		OracleScripts: oScripts,
		AIDataSources: aiDataSources,
		TestCases:     testCases,
		Params:        params,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		OracleScripts: []OracleScript{},
		AIDataSources: []AIDataSource{},
		TestCases:     []TestCase{},
		Params:        DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// ValidateGenesis validates the provider genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, record := range data.OracleScripts {
		if record.Name == "" {
			return fmt.Errorf("invalid OracleScripts: Value: %s. Error: Missing Name", record.Name)
		}
		if record.Owner == nil {
			return fmt.Errorf("invalid OracleScripts: Owner: %s. Error: Missing Owner", record.Owner)
		}
	}

	for _, record := range data.AIDataSources {
		if record.Name == "" {
			return fmt.Errorf("invalid AIDataSources: Value: %s. Error: Missing Name", record.Name)
		}
		if record.Owner == nil {
			return fmt.Errorf("invalid AIDataSources: Owner: %s. Error: Missing Owner", record.Owner)
		}
	}

	for _, record := range data.TestCases {
		if len(record.Owner) == 0 {
			return fmt.Errorf("invalid request: Owner: %s. Error: Missing requestID", string(record.Owner[:]))
		}
		if len(record.Name) == 0 {
			return fmt.Errorf("invalid request: Name: %s. Error: Missing name", record.Name)
		}
	}
	return nil
}
