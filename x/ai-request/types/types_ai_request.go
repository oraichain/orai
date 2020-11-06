package types

import (
	"fmt"
	"strings"

	"github.com/oraichain/orai/x/ai-request/exported"

	"github.com/oraichain/orai/x/provider/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// implements AIRequestI interface
var _ exported.AIRequestI = AIRequest{}

// AIRequest stores the request metadata to an AI model
type AIRequest struct {
	RequestID        string               `json:"request_id"`
	OracleScriptName string               `json:"oscript_name"`
	Creator          sdk.AccAddress       `json:"request_creator"`
	Validators       []sdk.ValAddress     `json:"validator_addr"`
	BlockHeight      int64                `json:"block_height"`
	AIDataSources    []types.AIDataSource `json:"data_sources"`
	TestCases        []types.TestCase     `json:"testcases"`
	Fees             sdk.Coins            `json:"transaction_fee"`
	Input            string               `json:"request_input"`
	ExpectedOutput   string               `json:"expected_output"`
	//Time         string `json:"time"`
}

// implement fmt.Stringer
func (ai AIRequest) String() string {
	valString := fmt.Sprintln(ai.Validators)
	dataSourceString := fmt.Sprintln(ai.AIDataSources)
	testCaseString := fmt.Sprintln(ai.TestCases)
	return strings.TrimSpace(fmt.Sprintf(`RequestID: %s
	OracleScript name: %s Creator: %s Validators: %s BlockHeight: %d AIDataSources: %s TestCases: %s Fees: %s`, ai.RequestID, ai.OracleScriptName, string(ai.Creator[:]), valString, ai.BlockHeight, dataSourceString, testCaseString, ai.Fees.String()))
}
