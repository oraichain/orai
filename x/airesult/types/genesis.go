package types

import (
	"fmt"
)

// GenesisState - all provider state that must be provided at genesis
type GenesisState struct {
	AIRequestResults []AIRequestResult `json:"ai_results"`
	Params           Params            `json:"params"`
	// TODO: Fill out what is needed by the module for genesis
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(aiResults []AIRequestResult, params Params) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		AIRequestResults: aiResults,
		Params:           params,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		AIRequestResults: []AIRequestResult{},
		Params:           DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// ValidateGenesis validates the provider genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, record := range data.AIRequestResults {
		if record.RequestID == "" {
			return fmt.Errorf("invalid AIRequestResults: Value: %s. Error: Missing RequestID", record.RequestID)
		}
		if len(record.Status) == 0 {
			return fmt.Errorf("invalid AIRequestResults: Status: %s. Error: Missing Status", record.Status)
		}
	}
	return nil
}
