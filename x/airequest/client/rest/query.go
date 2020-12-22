package rest

import (
	"fmt"
	"net/http"
	"io"
	"time"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	Timeout = 5
)

func registerWebsocketQueryRoutes(cliCtx context.CLIContext, r *mux.Router, in <-chan bool, out chan<- bool) {
	// TODO: Define your GET REST AIDataSources

	r.HandleFunc("/websocket/health", healthCheckHandler(in, out)).Methods("GET")

}

func pingChannel(ch chan<- bool) {    
    ch <- true
}

func healthCheckHandler(in <-chan bool, out chan<- bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// A very simple health check.
		w.Header().Set("Content-Type", "application/json")
		
		
		// send signal
		go pingChannel(out)		

		// timeout 10 seconds
		timer := time.NewTimer(Timeout * time.Second)
				
		select {
			
		// something may wrong event websocket routine alive
		case ret := <-in:
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, fmt.Sprintf(`{"alive": %v}`, ret))
		case <-timer.C:
			w.WriteHeader(http.StatusRequestTimeout)
			io.WriteString(w, `{"alive": false}`)				
		}
		

		// In the future we could report back on the status of our DB, or our cache
		// (e.g. Redis) by performing a simple PING, and include them in the response.
		// io.WriteString(w, `{"alive": true}`)
	}

}



func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// TODO: Define your GET REST AIDataSources

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/{%s}", storeName, restName),
		queryAIRequestHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireqs", storeName),
		queryAIRequestIDsHandlerFn(cliCtx),
	).Methods("GET")
}

func queryAIRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireq/%s", storeName, id), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryAIRequestIDsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireqs", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
