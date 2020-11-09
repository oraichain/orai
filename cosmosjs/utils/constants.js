module.exports = {
  ORAI_URL: process.env.URL,
  // [WARNING] This mnemonic is just for the demo purpose. DO NOT USE THIS MNEMONIC for your own wallet.
  MNEMONIC: process.env.MNEMONIC,
  CHAIN_ID: process.env.CHAIN_ID,
  // address of the account that uses to sign the transaction
  ACC_ADDRESS: process.env.ACC_ADDRESS,
  FINISHED_STATUS: "finished",
  INTERVAL: 5,
  MESSAGE_TYPE: {
    SET_PRICE_REQUEST: "airequest/SetPriceRequest"
  },
  DENOM: "orai"
}