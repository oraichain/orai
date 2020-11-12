const request = require('request');
const constants = require('../utils/constants');
const axios = require('axios');

const paths = {
  PRICE_REQ: "/airequest/aireq/pricereq",
  KYC_REQ: "/airequest/aireq/kycreq",
  CLASSIFICATION_REQ: "/airequest/aireq/clreq",
  GET_REQUEST_ID: "/txs/",
  GET_FULL_REQUEST: "/airesult/fullreq/",
  GET_MINIMUM_FEES: "/provider/min_fees/"
}

function getRequestId(payload, callback) {
  axios({
    method: "GET",
    url: constants.ORAI_URL + paths.GET_REQUEST_ID + payload
  })
    .then((response) => {
      callback(true, response, null);
    })
    .catch((error) => {
      callback(false, null, error);
    });
}

function getMinimumFees(payload, callback) {
  axios({
    method: "GET",
    url: constants.ORAI_URL + paths.GET_MINIMUM_FEES + payload.oScriptName + "?val_num=" + payload.valNum
  })
    .then((response) => {
      callback(true, response, null);
    })
    .catch((error) => {
      callback(false, null, error);
    });
}

function getFullRequest(payload, callback) {
  axios({
    method: "GET",
    url: constants.ORAI_URL + paths.GET_FULL_REQUEST + payload
  })
    .then((response) => {
      callback(true, response, null);
    })
    .catch((error) => {
      callback(false, null, error);
    });
}

function createFormUnsignedTx(payload, path, callback) {
  axios({
    method: "POST",
    url: constants.ORAI_URL + path,
    data: payload,
    headers: payload.getHeaders()
  })
    .then((response) => {
      callback(true, response, null);
    })
    .catch((error) => {
      callback(false, null, error);
    });
}

function createUnsignedTx(payload, path, callback) {
  axios({
    method: "POST",
    url: constants.ORAI_URL + path,
    data: payload,
  })
    .then((response) => {
      callback(true, response, null);
    })
    .catch((error) => {
      callback(false, null, error);
    });
}

const api = {
  createUnsignedTx,
  getRequestId,
  getFullRequest,
  getMinimumFees,
  paths,
  createFormUnsignedTx
}

module.exports = api;