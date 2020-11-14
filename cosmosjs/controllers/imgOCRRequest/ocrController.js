const { validationResult } = require("express-validator/check");
const cosmosjs = require("@cosmostation/cosmosjs");
const request = require('request');
const constants = require('../../utils/constants')
const api = require('../../api/api');
const formData = require("./formData");

module.exports = {
  getOCR: (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(422).json({ errors: errors.array() });
    } else {
      console.log("request file: ", req.file)
      const oScriptName = req.body.oscript_name
      // check minimum fees provided by the user
      const payload = {
        oScriptName: oScriptName,
        valNum: req.body.validator_count
      }
      api.getMinimumFees(payload, (isSuccess, response, error) => {
        if (isSuccess) {
          // if the fees smaller than the minium required fees then we return error
          if (req.body.fees - response.data.result.minimum_fees <= 0) {
            return res.status(200).json({
              message: "The given fee is smaller / equal to the required fee or the oracle script does not exist. Provided fees should at least be higher",
              required_fee: response.data.result.minimum_fees + "orai"
            });
          } else {
            //collect the orai object that connects with a rest server of a full node and chain ID
            const orai = cosmosjs.network(constants.ORAI_URL, constants.CHAIN_ID);

            // get pair private public key to sign the tx
            const ecpairPriv = orai.getECPairPriv(constants.MNEMONIC);

            // address of the account that uses to sign the transaction
            const accAddress = constants.ACC_ADDRESS

            // input encoded in base64 for the oracle script
            const input = req.body.input

            // expected output from the user for the test case
            //const expectedOutput = "NTAwMA=="
            const expectedOutput = req.body.expected_output

            // fees paid for the transaction
            const fees = req.body.fees + constants.DENOM

            // collect orai account to get acc sequence and number
            orai.getAccounts(accAddress).then(acc => {
              // create an unsigned tx
              var form = formData.generateOCRData({
                "from": accAddress,
                "chain_id": constants.CHAIN_ID,
                "oracle_script_name": oScriptName,
                "input": input,
                "expected_output": expectedOutput,
                "fees": fees,
                //"img_path": req.file.path,
                "validator_count": req.body.validator_count.toString()
              })

              // push image onto ipfs
              ipfsForm = formData.generateIPFSData(req.file.path)
              api.addIPFS(ipfsForm, (isSuccess, response, error) => {
                if (isSuccess) {
                  // add image information from ipfs to send to cosmos
                  form.append("image_hash", response.data.Hash)
                  form.append("image_name", response.data.Name)
                  api.createFormUnsignedTx(form, api.paths.OCR_REQ, (isSuccess, response, error) => {
                    if (isSuccess) {
                      console.log("unsigned tx: ", response.data)
                      unsignedTx = response.data;

                      // create a new stdSignMsg for signing
                      let stdSignMsg = orai.newStdMsg({
                        msgs: unsignedTx.value.msg,
                        chain_id: constants.CHAIN_ID,
                        fee: unsignedTx.value.fee,
                        memo: unsignedTx.value.memo,
                        account_number: String(acc.result.value.account_number),
                        sequence: String(acc.result.value.sequence)
                      });

                      const signedTx = orai.sign(stdSignMsg, ecpairPriv);
                      //broadcast the tx to the network and receive response
                      orai.broadcast(signedTx).then(response => {
                        // wait about 5 seconds to query the tx hash. Count check if interval 5 times still cannot get response then stop
                        console.log("response: ", response)
                        // TODO: at this time, need to verify if the txhash is a success or not. Normally, if success => no "code" field, if fail then have "code" field when querying the transaction hash
                        let txHash = response.txhash
                        let count = 0
                        let counter = 0
                        let resps = []
                        let msgs = []
                        getReqIdInterval = setInterval(() => {
                          count++;
                          api.getRequestId(txHash, (isSuccess, response, error) => {
                            if (isSuccess) {
                              // set to this value to stop the interval
                              count = constants.INTERVAL;
                              msgs = response.data.tx.value.msg
                              // loop through the list of msgs to collect request ids
                              for (let i = 0; i < msgs.length; i++) {
                                // since there can be more than one msg with different types, we need to do a msg type check. We only handle SET PRICE REQUEST type
                                if (msgs[i].type === constants.MESSAGE_TYPE.SET_OCR_REQUEST) {
                                  let resp = {
                                    requestId: [],
                                    validatorAddrs: [],
                                    blockHeight: "",
                                    aggregatedPrices: [],
                                    requestStatus: ""
                                  }
                                  resp.requestId = msgs[i].value.msg_set_ai_request.request_id
                                  // get full req information
                                  getFullReqInterval = setInterval(() => {
                                    counter++;
                                    api.getFullRequest(resp.requestId, (isSuccess, response, error) => {
                                      if (isSuccess) {
                                        console.log("status: ", response.data.result.ai_result.request_status)
                                        if (response.data.result.ai_result.request_status === constants.FINISHED_STATUS) {
                                          // set to this value to stop the interval
                                          counter = constants.INTERVAL;
                                          resp.validatorAddrs = response.data.result.ai_request.validator_addr
                                          resp.blockHeight = response.data.result.ai_request.block_height
                                          resp.aggregatedPrices = response.data.result.ai_result.results
                                          resp.requestStatus = response.data.result.ai_result.request_status
                                          resps.push(resp)
                                        }
                                        console.log("counter: ", counter)
                                        // Wait the full request for total 3 * 5 = 15 secs. If request still pending then we move on to the next loop
                                        if (counter >= constants.INTERVAL) {
                                          console.log("clear interval in get full request")
                                          clearInterval(getFullReqInterval);
                                          // if we already loop through all of the msgs then we return response
                                          if (i === msgs.length - 1) {
                                            // case where at least one requestId is has finished status
                                            if (resps.length > 0) {
                                              return res.status(200).json({
                                                message: "You have successfully signed and broadcast your tx",
                                                tx_hash: txHash,
                                                results: resps,
                                              });
                                            } else {
                                              // if we cannot obtain a valid request id (request still pending)
                                              return res.status(200).json({
                                                message: "The requests are still being handled on Oraichain. Please check them later using the transaction hash",
                                                tx_hash: txHash,
                                              });
                                            }
                                          }
                                        }
                                      } else {
                                        if (counter >= constants.INTERVAL) {
                                          clearInterval(getFullReqInterval);
                                          return res.status(400).json({
                                            message: "There is an error while getting the full request given a request id: " + resp.requestId,
                                            error: error,
                                          })
                                        }
                                      }
                                    })
                                  }, 3000)
                                } else {
                                  console.log("Skip request with message type: ", msgs[i].type)
                                }
                              }
                              // If we got the information from the tx hash, we stop looping the fetching request id
                              if (count >= constants.INTERVAL) {
                                clearInterval(getReqIdInterval);
                              }
                            } else {
                              // If after 5 times, we still cannot get information from the tx hash, we return error
                              if (count >= constants.INTERVAL) {
                                clearInterval(getReqIdInterval);
                                return res.status(400).json({
                                  message: "There is an error while getting request id. Please try again sometimes using the transaction hash: " + txHash,
                                  error: error
                                })
                              }

                            }
                          })
                        }, 3000);
                      });
                    } else {
                      return res.status(400).json({
                        message: "There is an error while creating an unsigned tx: ",
                        error: error,
                      })
                    }
                  })
                }
                else {
                  return res.status(400).json({
                    message: "There is an error while pushing image to IPFS: ",
                    error: error,
                  })
                }
              })
            })
          }
        } else {
          return res.status(404).json({
            message: "Cannot get the minimum fees of the oracle script",
            error: error
          });
        }
      })
    }
  },

  getOCRHash: (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(422).json({ errors: errors.array() });
    } else {
      console.log("request file: ", req.file)
      const oScriptName = req.body.oscript_name
      // check minimum fees provided by the user
      const payload = {
        oScriptName: oScriptName,
        valNum: req.body.validator_count
      }
      api.getMinimumFees(payload, (isSuccess, response, error) => {
        if (isSuccess) {
          // if the fees smaller than the minium required fees then we return error
          if (req.body.fees - response.data.result.minimum_fees <= 0) {
            return res.status(200).json({
              message: "The given fee is smaller / equal to the required fee or the oracle script does not exist. Provided fees should at least be higher",
              required_fee: response.data.result.minimum_fees + "orai"
            });
          } else {
            //collect the orai object that connects with a rest server of a full node and chain ID
            const orai = cosmosjs.network(constants.ORAI_URL, constants.CHAIN_ID);

            // get pair private public key to sign the tx
            const ecpairPriv = orai.getECPairPriv(constants.MNEMONIC);

            // address of the account that uses to sign the transaction
            const accAddress = constants.ACC_ADDRESS

            // input encoded in base64 for the oracle script
            const input = req.body.input

            // expected output from the user for the test case
            //const expectedOutput = "NTAwMA=="
            const expectedOutput = req.body.expected_output

            // fees paid for the transaction
            const fees = req.body.fees + constants.DENOM

            // collect orai account to get acc sequence and number
            orai.getAccounts(accAddress).then(acc => {
              // create an unsigned tx
              var form = formData.generateOCRData({
                "from": accAddress,
                "chain_id": constants.CHAIN_ID,
                "oracle_script_name": oScriptName,
                "input": input,
                "expected_output": expectedOutput,
                "fees": fees,
                //"img_path": req.file.path,
                "validator_count": req.body.validator_count.toString()
              })

              // push image onto ipfs
              form.append("image_hash", req.body.image_hash)
              form.append("image_name", req.body.image_name)
              api.createFormUnsignedTx(form, api.paths.OCR_REQ, (isSuccess, response, error) => {
                if (isSuccess) {
                  console.log("unsigned tx: ", response.data)
                  unsignedTx = response.data;

                  // create a new stdSignMsg for signing
                  let stdSignMsg = orai.newStdMsg({
                    msgs: unsignedTx.value.msg,
                    chain_id: constants.CHAIN_ID,
                    fee: unsignedTx.value.fee,
                    memo: unsignedTx.value.memo,
                    account_number: String(acc.result.value.account_number),
                    sequence: String(acc.result.value.sequence)
                  });

                  const signedTx = orai.sign(stdSignMsg, ecpairPriv);
                  //broadcast the tx to the network and receive response
                  orai.broadcast(signedTx).then(response => {
                    // wait about 5 seconds to query the tx hash. Count check if interval 5 times still cannot get response then stop
                    console.log("response: ", response)
                    // TODO: at this time, need to verify if the txhash is a success or not. Normally, if success => no "code" field, if fail then have "code" field when querying the transaction hash
                    let txHash = response.txhash
                    let count = 0
                    let counter = 0
                    let resps = []
                    let msgs = []
                    getReqIdInterval = setInterval(() => {
                      count++;
                      api.getRequestId(txHash, (isSuccess, response, error) => {
                        if (isSuccess) {
                          // set to this value to stop the interval
                          count = constants.INTERVAL;
                          msgs = response.data.tx.value.msg
                          // loop through the list of msgs to collect request ids
                          for (let i = 0; i < msgs.length; i++) {
                            // since there can be more than one msg with different types, we need to do a msg type check. We only handle SET PRICE REQUEST type
                            if (msgs[i].type === constants.MESSAGE_TYPE.SET_OCR_REQUEST) {
                              let resp = {
                                requestId: [],
                                validatorAddrs: [],
                                blockHeight: "",
                                aggregatedPrices: [],
                                requestStatus: ""
                              }
                              resp.requestId = msgs[i].value.msg_set_ai_request.request_id
                              // get full req information
                              getFullReqInterval = setInterval(() => {
                                counter++;
                                api.getFullRequest(resp.requestId, (isSuccess, response, error) => {
                                  if (isSuccess) {
                                    console.log("status: ", response.data.result.ai_result.request_status)
                                    if (response.data.result.ai_result.request_status === constants.FINISHED_STATUS) {
                                      // set to this value to stop the interval
                                      counter = constants.INTERVAL;
                                      resp.validatorAddrs = response.data.result.ai_request.validator_addr
                                      resp.blockHeight = response.data.result.ai_request.block_height
                                      resp.aggregatedPrices = response.data.result.ai_result.results
                                      resp.requestStatus = response.data.result.ai_result.request_status
                                      resps.push(resp)
                                    }
                                    console.log("counter: ", counter)
                                    // Wait the full request for total 3 * 5 = 15 secs. If request still pending then we move on to the next loop
                                    if (counter >= constants.INTERVAL) {
                                      console.log("clear interval in get full request")
                                      clearInterval(getFullReqInterval);
                                      // if we already loop through all of the msgs then we return response
                                      if (i === msgs.length - 1) {
                                        // case where at least one requestId is has finished status
                                        if (resps.length > 0) {
                                          return res.status(200).json({
                                            message: "You have successfully signed and broadcast your tx",
                                            tx_hash: txHash,
                                            results: resps,
                                          });
                                        } else {
                                          // if we cannot obtain a valid request id (request still pending)
                                          return res.status(200).json({
                                            message: "The requests are still being handled on Oraichain. Please check them later using the transaction hash",
                                            tx_hash: txHash,
                                          });
                                        }
                                      }
                                    }
                                  } else {
                                    if (counter >= constants.INTERVAL) {
                                      clearInterval(getFullReqInterval);
                                      return res.status(400).json({
                                        message: "There is an error while getting the full request given a request id: " + resp.requestId,
                                        error: error,
                                      })
                                    }
                                  }
                                })
                              }, 3000)
                            } else {
                              console.log("Skip request with message type: ", msgs[i].type)
                            }
                          }
                          // If we got the information from the tx hash, we stop looping the fetching request id
                          if (count >= constants.INTERVAL) {
                            clearInterval(getReqIdInterval);
                          }
                        } else {
                          // If after 5 times, we still cannot get information from the tx hash, we return error
                          if (count >= constants.INTERVAL) {
                            clearInterval(getReqIdInterval);
                            return res.status(400).json({
                              message: "There is an error while getting request id. Please try again sometimes using the transaction hash: " + txHash,
                              error: error
                            })
                          }

                        }
                      })
                    }, 3000);
                  });
                } else {
                  return res.status(400).json({
                    message: "There is an error while creating an unsigned tx: ",
                    error: error,
                  })
                }
              })
            })
          }
        } else {
          return res.status(404).json({
            message: "Cannot get the minimum fees of the oracle script",
            error: error
          });
        }
      })
    }
  },
}
