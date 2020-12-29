package rest

import (
	"bytes"
	goContext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	dockerType "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gorilla/mux"
	"github.com/oraichain/orai/x/provider/types"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// TODO: Define your GET REST AIDataSources
	r.HandleFunc(
		fmt.Sprintf("/%s/oscript/{%s}", storeName, restName),
		queryOracleScriptHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/onames", storeName),
		queryOracleScriptNamesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/oscripts", storeName),
		queryOracleScriptsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/datasource/{%s}", storeName, restName),
		queryDataSourceHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/datasources", storeName),
		queryDataSourcesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/dnames", storeName),
		queryDataSourceNamesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/testcase/{%s}", storeName, restName),
		queryTestCaseHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/testcases", storeName),
		queryTestCasesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/tcnames", storeName),
		queryTestCaseNamesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/min_fees/{%s}", storeName, restName),
		queryMinimumFeesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/spec_oscript/{%s}", storeName, restName),
		querySpecPythonFile(cliCtx),
	).Methods("GET")
}

func queryOracleScriptHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/oscript/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryOracleScriptsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get query params
		v := r.URL.Query()
		page := v.Get("page")
		limit := v.Get("limit")
		name := v.Get("name")

		// In case the request does not include pagination parameters
		if page == "" || limit == "" {
			page = types.DefaultQueryPage
			limit = types.DefaultQueryLimit
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/oscripts", storeName), []byte(page+"-"+limit+"-"+name))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryOracleScriptNamesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/onames", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDataSourceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/datasource/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDataSourcesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		page := v.Get("page")
		limit := v.Get("limit")
		name := v.Get("name")

		// In case the request does not include pagination parameters
		if page == "" || limit == "" {
			page = types.DefaultQueryPage
			limit = types.DefaultQueryLimit
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/datasources", storeName), []byte(page+"-"+limit+"-"+name))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDataSourceNamesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/dnames", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTestCaseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/testcase/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTestCasesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		page := v.Get("page")
		limit := v.Get("limit")
		name := v.Get("name")

		// In case the request does not include pagination parameters
		if page == "" || limit == "" {
			page = types.DefaultQueryPage
			limit = types.DefaultQueryLimit
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/testcases", storeName), []byte(page+"-"+limit+"-"+name))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTestCaseNamesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tcnames", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryMinimumFeesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]
		v := r.URL.Query()
		valNum := v.Get("val_num")

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/min_fees/%s", storeName, name), []byte(valNum))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySpecPythonFile(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]
		ctx := goContext.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			rest.PostProcessResponse(w, cliCtx, err.Error())
		}
		script := "" +
			"import yaml, json\n" +
			"import sys\n" +
			"\n" +
			"if __name__ == \"__main__\":\n" +
			"    try:\n" +
			"        data_source = __import__(sys.argv[1])\n" +
			"        doc = data_source.__doc__\n" +
			"        if doc:\n" +
			"            comment_index = doc.rfind('---')\n" +
			"            if comment_index > 0:\n" +
			"                comment_index = comment_index + 3\n" +
			"            else:\n" +
			"                comment_index = 0\n" +
			"            code = yaml.safe_load((doc[comment_index:]))\n" +
			"            print(json.dumps(code))\n" +
			"    except ArithmeticError:\n" +
			"        print(\"\")"
		resp, err := cli.ContainerExecCreate(ctx, "python", dockerType.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			Cmd:          append([]string{"python", "-c", script, name}),
		})
		if err != nil {
			rest.PostProcessResponse(w, cliCtx, err.Error())
		}

		logResp, err := cli.ContainerExecAttach(ctx, resp.ID, dockerType.ExecStartCheck{})
		if err != nil {
			rest.PostProcessResponse(w, cliCtx, err.Error())
		}

		var buf, error bytes.Buffer
		stdcopy.StdCopy(&buf, &error, logResp.Reader)
		rest.PostProcessResponse(w, cliCtx, string(buf.Bytes()))
	}
}
