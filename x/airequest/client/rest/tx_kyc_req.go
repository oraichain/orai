package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/segmentio/ksuid"
)

type setKYCRequestReq struct {
	FormAIRequest setFormAIRequestReq `json:"base_form_request"`
	Name          string              `json:"Name"`
	Hash          string              `json:"Hash"`
}

// newSetKYCRequestReq is the constructor for the setKYCRequestReq
func newSetKYCRequestReq(formAIReq setFormAIRequestReq, name, hash string) setKYCRequestReq {
	return setKYCRequestReq{
		FormAIRequest: formAIReq,
		Name:          name,
		Hash:          hash,
	}
}

func (kyc setKYCRequestReq) getFormAIRequest() setFormAIRequestReq {
	return kyc.FormAIRequest
}

func (kyc setKYCRequestReq) getName() string {
	return kyc.Name
}

func (kyc setKYCRequestReq) getHash() string {
	return kyc.Hash
}

// setKYCRequestReqFn is the function that collects all the necessary info of KYC and return a new object out of it
func setKYCRequestReqFn(cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request) setKYCRequestReq {
	req := setAIRequestHandlerFn(cliCtx, w, r)
	imageHash := r.FormValue("image_hash")
	imageName := r.FormValue("image_name")
	return newSetKYCRequestReq(req, imageHash, imageName)
}

func setKYCRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// // collect image file from user
		// file, handler, err := r.FormFile("image")
		// if err != nil {
		// 	fmt.Println("Error Retrieving the File")
		// 	fmt.Println(err)
		// 	return
		// }
		// defer file.Close()

		// fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")

		// // Create a temp file in local storage for IPFS http request
		// tempFile, err := ioutil.TempFile("./", "upload-*.png")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// defer tempFile.Close()
		// fileBytes, err := ioutil.ReadAll(file)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// tempFile.Write(fileBytes)

		// fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

		// // Prepare to send the image onto IPFS
		// b, writer, err := filehandling.CreateMultipartFormData("image", tempFile.Name())

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to create multiform data image")
		// 	return
		// }

		// httpReq, err := http.NewRequest("POST", types.IPFSUrl+types.IPFSAdd, &b)
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to create new request for sending image to IPFS")
		// 	return
		// }
		// // Don't forget to set the content type, this will contain the boundary.
		// httpReq.Header.Set("Content-Type", writer.FormDataContentType())

		// client := &http.Client{}
		// resp, err := client.Do(httpReq)

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to execute request for sending image to IPFS")
		// 	return
		// }

		// defer resp.Body.Close()

		// result := setIPFSImage{}

		// // Collect the result in json form. Remember that we need to create a corresponding struct to do this
		// json.NewDecoder(resp.Body).Decode(&result)

		// // After collecting the hash image, we need to clear the image file stored temporary
		// err = os.Remove(tempFile.Name())
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to remove the temporary image file: %s", err.Error()))
		// 	return
		// }

		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to push image onto IPFS: %s", err.Error()))
		// 	return
		// }

		req := setKYCRequestReqFn(cliCtx, w, r)

		// Need to create a baseReq to write tx response. We cannot use baseReq in the AIRequest struct because AIRequest needs to be in form data to be able to send images
		baseReq := rest.BaseReq{
			From:          req.getFormAIRequest().From,
			Memo:          req.getFormAIRequest().Memo,
			ChainID:       req.getFormAIRequest().ChainID,
			AccountNumber: req.getFormAIRequest().AccountNumber,
			Sequence:      req.getFormAIRequest().Sequence,
			Fees:          req.getFormAIRequest().Fees,
			GasPrices:     req.getFormAIRequest().GasPrices,
			Gas:           req.getFormAIRequest().Gas,
			GasAdjustment: req.getFormAIRequest().GasAdjustment,
			Simulate:      req.getFormAIRequest().Simulate,
		}

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(req.getFormAIRequest().From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "AVXSD")
			return
		}

		// create the message
		// msg := types.NewMsgSetKYCRequest(result.Hash, handler.Filename, types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees.String(), req.ValidatorCount, req.Input, req.ExpectedOutput))
		msg := types.NewMsgSetKYCRequest(req.getHash(), req.getName(), types.NewMsgSetAIRequest(ksuid.New().String(), req.getFormAIRequest().OracleScriptName, addr, req.getFormAIRequest().Fees.String(), req.getFormAIRequest().ValidatorCount, req.getFormAIRequest().Input, req.getFormAIRequest().ExpectedOutput))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GHYK")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
