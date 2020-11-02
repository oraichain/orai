const request = require('request');
const constants = require('../utils/constants');
const axios = require('axios');

const paths = {
  CREATE_UNSIGNED_TX: "/provider/aireq/pricereq",
  GET_REQUEST_ID: "/txs/",
  GET_FULL_REQUEST: "/provider/fullreq/",
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

function createUnsignedTx(payload, callback) {
  axios({
    method: "POST",
    url: constants.ORAI_URL + paths.CREATE_UNSIGNED_TX,
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
  getMinimumFees
}

module.exports = api;