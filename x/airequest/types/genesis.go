package types

import (
	"fmt"
)

// GenesisState - all provider state that must be provided at genesis
type GenesisState struct {
	AIRequests []AIRequest `json:"ai_requests"`
	Params     Params      `json:"params"`
	// TODO: Fill out what is needed by the module for genesis
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(aiRequests []AIRequest, params Params) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		AIRequests: aiRequests,
		Params:     params,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		AIRequests: []AIRequest{},
		Params:     DefaultParams(),

		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	}
}

// ValidateGenesis validates the provider genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, record := range data.AIRequests {
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
