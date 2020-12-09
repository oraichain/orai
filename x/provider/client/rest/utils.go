package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// getFileBytes try to parse the user path input and collect the raw bytes of the file.
func getFileBytes(path string) ([]byte, error) {
	// placeholder for the source code in bytes
	var execBytes []byte

	// validate the code path of the script
	u, err := url.ParseRequestURI(path)
	if err != nil {
		return nil, err
	}

	// check again if the path is a URL (http) or not
	u, err = url.Parse(path)
	if err != nil || u.Scheme == "" || u.Host == "" {
		// if error then it must be path local
		// collect the byte code of the source code based on the path
		fmt.Println("reading from local")
		execBytes, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return execBytes, nil
	}

	// fetch the file using http if no error
	fmt.Println("path is http: ", path)
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	// prepare to close the response body data after fetching.
	defer resp.Body.Close()

	// read the file data to collect the bytes version of the file
	execBytes, err = ioutil.ReadAll(resp.Body)
	return execBytes, nil
}
