var msgpack = require("msgpack-lite");

// encode from JS Object to MessagePack (Buffer)
var buffer = msgpack.encode({"foo": "bar"});

console.log("buffer: ", buffer)

// decode from MessagePack (Buffer) to JS Object
var data = msgpack.decode(buffer); // => {"foo": "bar"}

// if encode/decode receives an invalid argument an error is thrown

var axios = require('axios');
var FormData = require('form-data');
//var canvas = require('canvas');
var fs = require('fs');
var data = new FormData();
//data.append('oscript_name', 'oscript_classification');
data.append('image', fs.createReadStream('/home/duc/bkc_research/cosmos/oraichain/images/126185086_186028783118713_7125389096621276112_n.jpg'));
//data.append('input', 'Z29sZGVuX3JldHJpZXZlcg==');
//data.append('expected_output', 'Z29sZGVuX3JldHJpZXZlcg==');
//data.append('fees', '81000');
//data.append('validator_count', '1');
//var data = "abcdef"

var input = msgpack.encode("")
var expected_output = msgpack.encode({"image": data, "test": "gg"})
var obj = {"from":"orai1anmut6685agd7kxtjssawpd8u4mgak374mdzhy","chain_id":"Oraichain","oracle_script_name":"oscript_price","input":input,"expected_output":expected_output,"fees":"35000orai","validator_count":1};
var buf = msgpack.encode(obj)
console.log("buffer: ", buf)

var config = {
    method: 'post',
    url: 'http://localhost:1317/airequest/aireq/pricereq',
    headers: {
      'Content-Type': 'application/msgpack'
    },
    data: buf
  };
  
  axios(config)
  .then(function (response) {
    console.log(JSON.stringify(response.data));
  })
  .catch(function (error) {
    console.log(error);
  });
