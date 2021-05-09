package types

// provider module event types
const (
	// TODO: Create your event types
	// EventType<Action>    		= "action"
	// Event type for set ai oracle
	EventTypeRequestWithData        = "ai_oracle_data"
	EventTypeReportWithData         = "report_data"
	EventTypeTestCaseReportWithData = "testcase_report_data"
	EventTypeSetAIOracle            = "set_ai_oracle"
	EventTypeSetReport              = "set_report"
	EventTypeSetTestCaseReport      = "set_testcase_report"

	// TODO: Create keys fo your events, the values will be derivided from the msg
	// AttributeKeyAddress  		= "address"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
	// Attribute for request
	AttributeRequest        = "request"
	AttributeBaseReport     = "base_report"
	AttributeReport         = "report"
	AttributeTestCaseReport = "testcase_report"
	AttributeKeyValidator   = "validator"
)
