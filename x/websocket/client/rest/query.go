package rest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

const (
	timeout = 8 * time.Second
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, in <-chan bool, out chan<- bool) {
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

		// timeout 8 seconds
		timer := time.NewTimer(timeout)

		select {

		// something may wrong event websocket routine alive
		case ret := <-in:
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, fmt.Sprintf(`{"alive": %v}`, ret))
		case <-timer.C:
			w.WriteHeader(http.StatusRequestTimeout)
			io.WriteString(w, fmt.Sprintln(`{"alive": false}`))
		}

		// In the future we could report back on the status of our DB, or our cache
		// (e.g. Redis) by performing a simple PING, and include them in the response.
		// io.WriteString(w, `{"alive": true}`)
	}

}
