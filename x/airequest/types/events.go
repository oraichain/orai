package types

// provider module event types
const (
	// TODO: Create your event types
	// EventType<Action>    		= "action"
	// Event type for set ai request
	EventTypeRequestWithData = "ai_request_data"
	EventTypeSetAIRequest    = "set_ai_request"

	// TODO: Create keys fo your events, the values will be derivided from the msg
	// AttributeKeyAddress  		= "address"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
	// Attribute for request
	AttributeRequestID        = "request_id"
	AttributeRequestValidator = "request_validator"
	AttributeRequestCreator   = "request_creator"
	// Attributes for oracle script
	AttributeOracleScriptName      = "oscript_name"
	AttributeRequestImageHash      = "request_image_hash"
	AttributeRequestImageName      = "request_image_name"
	AttributeRequestValidatorCount = "request_validator_count"
	AttributeRequestInput          = "request_input"
	AttributeRequestExpectedOutput = "expected_request_output"
	AttributeRequestDSources       = "data_sources"
	AttributeRequestTCases         = "test_cases"
)
