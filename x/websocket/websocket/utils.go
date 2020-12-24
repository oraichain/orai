package websocket

import (
	"log"
	"os"
	"path"
	"regexp"

	"github.com/oraichain/orai/x/provider"
)

func getCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}

func getOScriptPath(name string) string {
	return path.Join(getCurrentDir(), provider.ScriptPath, provider.OracleScriptFileString(name))
}

func getDSourcePath(name string) string {
	return path.Join(getCurrentDir(), provider.ScriptPath, provider.DataSourceStoreFileString(name))
}

func getTCasePath(name string) string {
	// get absolute path from working dir
	return path.Join(getCurrentDir(), provider.ScriptPath, provider.TestCaseStoreFileString(name))
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
