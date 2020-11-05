package types

// import (
// 	"fmt"
// 	"strings"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// // OracleScript is a struct for storing oracle script information
// type OracleScript struct {
// 	Name        string         `json:"name"`
// 	Owner       sdk.AccAddress `json:"owner"`
// 	Description string         `json:"description"`
// 	MinimumFees sdk.Coins      `json:"minimum_fees"`
// }

// // implement fmt.Stringer
// func (c OracleScript) String() string {
// 	return strings.TrimSpace(fmt.Sprintf(`Name: %s
// Owner: %s Description: %s Minimum Fees: %s`, c.Name, c.Owner, c.Description, c.MinimumFees))
// }

// // AIDataSource is a struct for storing AIDataSource information of a provider
// type AIDataSource struct {
// 	Name        string         `json:"name"`
// 	Owner       sdk.AccAddress `json:"owner"`
// 	Fees        sdk.Coins      `json:"transaction_fee"`
// 	Description string         `json:"description"`
// }

// // implement fmt.Stringer
// func (e AIDataSource) String() string {
// 	return strings.TrimSpace(fmt.Sprintf(`Name: %s
// Owner: %s Fees: %s Description: %s`, e.Name, e.Owner, e.Fees.String(), e.Description))
// }

// // AIRequest stores the request metadata to an AI model
// type AIRequest struct {
// 	RequestID        string           `json:"request_id"`
// 	OracleScriptName string           `json:"oscript_name"`
// 	Creator          sdk.AccAddress   `json:"request_creator"`
// 	Validators       []sdk.ValAddress `json:"validator_addr"`
// 	BlockHeight      int64            `json:"block_height"`
// 	AIDataSources    []AIDataSource   `json:"data_sources"`
// 	TestCases        []TestCase       `json:"testcases"`
// 	Fees             sdk.Coins        `json:"transaction_fee"`
// 	Input            string           `json:"request_input"`
// 	ExpectedOutput   string           `json:"expected_output"`
// 	//Time         string `json:"time"`
// }

// // implement fmt.Stringer
// func (ai AIRequest) String() string {
// 	valString := fmt.Sprintln(ai.Validators)
// 	dataSourceString := fmt.Sprintln(ai.AIDataSources)
// 	testCaseString := fmt.Sprintln(ai.TestCases)
// 	return strings.TrimSpace(fmt.Sprintf(`RequestID: %s
// 	OracleScript name: %s Creator: %s Validators: %s BlockHeight: %d AIDataSources: %s TestCases: %s Fees: %s`, ai.RequestID, ai.OracleScriptName, string(ai.Creator[:]), valString, ai.BlockHeight, dataSourceString, testCaseString, ai.Fees.String()))
// }

// // TestCase stores the test case of a request to an AI model
// type TestCase struct {
// 	Owner       sdk.AccAddress `json:"owner"`
// 	Name        string         `json:"name"`
// 	Fees        sdk.Coins      `json:"transaction_fee"`
// 	Description string         `json:"description"`
// }

// // implement fmt.Stringer
// func (tc TestCase) String() string {
// 	return strings.TrimSpace(fmt.Sprintf(`
// Owner: %s Name: %s Fees: %s Description: %s`, tc.Owner, tc.Name, tc.Fees.String(), tc.Description))
// }

// // Validator mimics the original validator to store information of those that execute the oScript
// type Validator struct {
// 	Address     sdk.ValAddress `json:"address"`
// 	VotingPower int64          `json:"voting_power"`
// 	Status      string         `json:"status"`
// }

// // DataSourceResult stores the data source result
// type DataSourceResult struct {
// 	Name   string `json:"data_source"`
// 	Result []byte `json:"result"`
// 	Status string `json:"result_status"`
// }

// // TestCaseResult stores the test case result
// type TestCaseResult struct {
// 	Name              string             `json:"test_case"`
// 	DataSourceResults []DataSourceResult `json:"data_source_result"`
// }

// // Report stores the result of the data source when validator executes it
// type Report struct {
// 	RequestID         string             `json:"request_id"`
// 	Validator         Validator          `json:"validator"`
// 	DataSourceResults []DataSourceResult `json:"data_source_results"`
// 	TestCaseResults   []TestCaseResult   `json:"test_case_results"`
// 	BlockHeight       int64              `json:"block_height"`
// 	Fees              sdk.Coins          `json:"report_fee"`
// 	AggregatedResult  []byte             `json:"aggregated_result"`
// }

// // Reward stores a list of validators, data source owners and test case owners that receive rewards for a specific block height
// type Reward struct {
// 	Validators    []Validator    `json:"validators"`
// 	DataSources   []AIDataSource `json:"data_sources"`
// 	TestCases     []TestCase     `json:"test_cases"`
// 	BlockHeight   int64          `json:"block_height"`
// 	TotalPower    int64          `json:"total_voting_power"`
// 	ProviderFees  sdk.Coins      `json:"provider_fees"`
// 	ValidatorFees sdk.Coins      `json:"validator_fees"`
// }

// // AIRequestResult stores the final result after aggregating the results from the reports of an AI request
// type AIRequestResult struct {
// 	RequestID string     `json:"request_id"`
// 	Results   ValResults `json:"results"`
// 	Status    string     `json:"request_status"`
// }

// // ValResult stores the result information from a validator that has executed the oracle script
// type ValResult struct {
// 	Validator sdk.ValAddress `json:"validator_address"`
// 	Result    []byte         `json:"result"`
// }

// // ValResults is the list of results struct
// type ValResults []ValResult

// // Strategy stores the information of a strategy for a yAI flow
// type Strategy struct {
// 	StratID        uint64   `json:"strategy_id"`
// 	StratName      string   `json:"strategy_name"`
// 	StratFlow      []string `json:"strategy_flow"`
// 	PerformanceFee uint64   `json:"performance_fee"`
// 	PerformanceMax uint64   `json:"performance_max"`
// 	WithdrawalFee  uint64   `json:"withdrawal_fee"`
// 	WithdrawalMax  uint64   `json:"withdrawal_max"`
// 	GovernanceAddr string   `json:"governance_address"`
// 	StrategistAddr string   `json:"strategist_address"`
// }

// // implement fmt.Stringer
// func (s Strategy) String() string {
// 	stratFlows := fmt.Sprintln(s.StratFlow)
// 	return strings.TrimSpace(fmt.Sprintf(`StratID: %d
// 	StratName: %s StratFlow: %s PerformanceFee: %s PerformanceMax %d WithdrawalFee %d WithdrawalMax %d GovernanceAddr %s StrategistAddr %s`, s.StratID, s.StratName, stratFlows, s.PerformanceFee, s.PerformanceMax, s.WithdrawalFee, s.WithdrawalMax, s.GovernanceAddr, s.StrategistAddr))
// }
