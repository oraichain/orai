package types

// provider module event types
const (
	// TODO: Create your event types
	// EventType<Action>    		= "action"
	// Event type for oracle script
	EventTypeSetOracleScript  = "set_oscript"
	EventTypeEditOracleScript = "edit_oscript"
	// Event type for data source
	EventTypeSetDataSource  = "set_datasource"
	EventTypeEditDataSource = "edit_datasource"
	// Event type for test case
	EventTypeCreateTestCase = "set_test_case"
	EventTypeEditTestCase   = "set_edit_case"

	// TODO: Create keys fo your events, the values will be derivided from the msg
	// AttributeKeyAddress  		= "address"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
	// Attributes for oracle script
	AttributeOracleScriptName = "oscript_name"
	// Attribute for data source
	AttributeDataSourceName = "datasource_name"

	// Attribute for test case
	AttributeTestCaseName = "test_case_name"
)
