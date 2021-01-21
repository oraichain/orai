package websocket

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	aiRequest "github.com/oraichain/orai/x/airequest"
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

		msg, isProvider := verifyProviderMessageType(messageType, log)

		if messageType == (aiRequest.MsgSetAIRequest{}).Type() {
			go handleAIRequestLog(c, l, log)
		} else if isProvider {
			go handleProviderMsgLog(c, l, log, msg)
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

// getPaths collect data sources and test cases information from the ai request event
func getPaths(log sdk.ABCIMessageLog) ([]string, []string, error) {
	aiDataSourcesStr, err := GetEventValue(log, aiRequest.EventTypeRequestWithData, aiRequest.AttributeRequestDSources)
	if err != nil {
		return nil, nil, err
	}

	testCasesStr, err := GetEventValue(log, aiRequest.EventTypeRequestWithData, aiRequest.AttributeRequestTCases)
	if err != nil {
		return nil, nil, err
	}
	return strings.Split(aiDataSourcesStr, "-"), strings.Split(testCasesStr, "-"), nil
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
