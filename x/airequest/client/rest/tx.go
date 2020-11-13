package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (

	// "bytes"
	// "net/http"

	"fmt"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	// "github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/kycreq", storeName),
		setKYCRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/clreq", storeName),
		setClassificationRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/ocrreq", storeName),
		setOCRRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/pricereq", storeName),
		setPriceRequestHandlerFn(cliCtx),
	).Methods("POST")
}
