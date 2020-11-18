package websocket

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	aiRequest "github.com/oraichain/orai/x/airequest/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"
)

func handleTransaction(c *Context, l *Logger, tx tmtypes.TxResult) {
	l.Debug(":eyes: Inspecting incoming transaction: %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		l.Debug(":alien: Skipping transaction with non-zero code: %d", tx.Result.Code)
		return
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		l.Error(":cold_sweat: Failed to parse transaction logs with error: %s", err.Error())
		return
	}

	for _, log := range logs {
		messageType, err := GetEventValue(log, sdk.EventTypeMessage, sdk.AttributeKeyAction)
		if err != nil {
			l.Error(":cold_sweat: Failed to get message action type with error: %s", err.Error())
			continue
		}

		l.Info(":star: message type: %s", messageType)

		if messageType == (aiRequest.MsgSetKYCRequest{}).Type() {
			go handleKYCRequestLog(c, l, log)
		} else if messageType == (aiRequest.MsgSetPriceRequest{}).Type() {
			go handlePriceRequestLog(c, l, log)
		} else if messageType == (aiRequest.MsgSetClassificationRequest{}).Type() {
			go handleClassificationRequestLog(c, l, log)
		} else if messageType == (aiRequest.MsgSetOCRRequest{}).Type() {
			go handleOCRRequestLog(c, l, log)
		} else {
			l.Debug(":ghost: Skipping non-{request/packet} type: %s", messageType)
		} /*else if messageType == (ibc.MsgPacket{}).Type() {
			// Try to get request id from packet. If not then return error.
			_, err := GetEventValue(log, types.EventTypeRequest, types.AttributeKeyID)
			if err != nil {
				l.Debug(":ghost: Skipping non-request packet")
				return
			}
			go handleRequestLog(c, l, log)
		} */
	}
}

func getPaths(l *Logger, oscriptPath string) [][]string {
	//use "data source" as an argument to collect the data source script name
	cmd := exec.Command("bash", oscriptPath, "aiDataSource")
	var dataSourceName bytes.Buffer
	cmd.Stdout = &dataSourceName
	err := cmd.Run()
	if err != nil {
		l.Error(":skull: failed to collect data source name: %s", err.Error())
	}

	// collect data source result from the script
	result := strings.TrimSuffix(dataSourceName.String(), "\n")

	aiDataSources := strings.Fields(result)

	//use "test case" as an argument to collect the test case script name
	cmd = exec.Command("bash", oscriptPath, "testcase")
	var testCaseName bytes.Buffer
	cmd.Stdout = &testCaseName
	err = cmd.Run()
	if err != nil {
		l.Error(":skull: failed to collect test case name: %s", err.Error())
	}

	// collect data source result from the script
	result = strings.TrimSuffix(testCaseName.String(), "\n")

	testCases := strings.Fields(result)

	var listPaths [][]string
	listPaths = append(listPaths, aiDataSources)
	listPaths = append(listPaths, testCases)

	return listPaths
}

// GetEventValues returns the list of all values in the given log with the given type and key.
func GetEventValues(log sdk.ABCIMessageLog, evType string, evKey string) (res []string) {
	for _, ev := range log.Events {
		fmt.Println(":delivery_truck: event collected: ", ev.Type)
		if ev.Type != evType {
			continue
		}

		for _, attr := range ev.Attributes {
			fmt.Println("Attribute key: ", attr.Key)
			fmt.Println("Attribute value: ", attr.Value)
			if attr.Key == evKey {
				res = append(res, attr.Value)
			}
		}
	}
	return res
}

// GetEventValue checks and returns the exact value in the given log with the given type and key.
func GetEventValue(log sdk.ABCIMessageLog, evType string, evKey string) (string, error) {

	values := GetEventValues(log, evType, evKey)
	if len(values) == 0 {
		return "", fmt.Errorf("Cannot find event with type: %s, key: %s", evType, evKey)
	}
	if len(values) > 1 {
		return "", fmt.Errorf("Found more than one event with type: %s, key: %s", evType, evKey)
	}
	return values[0], nil
}
