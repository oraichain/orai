package websocket

import (
	"log"
	"os"
	"path"
	"regexp"

	"github.com/oraichain/orai/x/provider/types"
	provider "github.com/oraichain/orai/x/provider/types"
)

func getCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}

func getOScriptPath(name string) string {
	return path.Join(getCurrentDir(), provider.ScriptPath, types.OracleScriptStoreKeyString(name))
}

func getDSourcePath(name string) string {
	return path.Join(getCurrentDir(), provider.ScriptPath, types.DataSourceStoreFileString(name))
}

func getTCasePath(name string) string {
	// get absolute path from working dir
	return path.Join(getCurrentDir(), provider.ScriptPath, types.TestCaseStoreFileString(name))
}

// trimResultEscapeChars is used to remove all escape characters in a string for cleaning the result
func trimResultEscapeChars(result string) string {
	re, err := regexp.Compile(`[\n\t\r]`)
	if err != nil {
		log.Fatal(err)
	}
	res := re.ReplaceAllString(result, "")
	return res
}

const (
	// ScriptPath is the path that stores all the script files (oracle scripts, data sources, test cases)
	ScriptPath = "../../../.oraifiles/"
	// Delimiter is the delimiter for separating results
	delimiter = "-"
)
