package rest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	// "github.com/oraichain/orai/x/websocket/websocket"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, in <-chan bool, out chan<- bool) {
	// TODO: Define your GET REST AIDataSources

	r.HandleFunc("/websocket/health", healthCheckHandler(in, out)).Methods("GET")

}

func healthCheckHandler(in <-chan bool, out chan<- bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// A very simple health check.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// send signal
		out <- true

		// timeout 10 seconds
		select {
		// something may wrong event websocket routine alive
		case ret := <-in:
			io.WriteString(w, fmt.Sprintf(`{"alive": %v}`, ret))
		case <-time.After(10 * time.Second):
			io.WriteString(w, `{"alive": false}`)
		}

		// In the future we could report back on the status of our DB, or our cache
		// (e.g. Redis) by performing a simple PING, and include them in the response.
		// io.WriteString(w, `{"alive": true}`)
	}

}
