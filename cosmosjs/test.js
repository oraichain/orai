var msgpack = require("msgpack-lite");

// if encode/decode receives an invalid argument an error is thrown

var axios = require('axios');
//var canvas = require('canvas');
var fs = require('fs');

fs.readFile('/home/duc/bkc_research/cosmos/oraichain/images/126185086_186028783118713_7125389096621276112_n.jpg', function (err, data) {
  if (err) throw err;
  console.log("image data: ", data);

  var input = msgpack.encode("")
  var expected_output = msgpack.encode({ "image": "dataxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxfzxfzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzsafasfasfsaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaawfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffsxxxxxxxxxxxxxxxxx", "test": "gg", "x": "y", "Z": "g" })
  var obj = { "from": "orai1anmut6685agd7kxtjssawpd8u4mgak374mdzhy", "chain_id": "Oraichain", "oracle_script_name": "oscript_price", "input": input, "expected_output": expected_output, "fees": "35000orai", "validator_count": 1 };
  var buf = msgpack.encode(obj)
  console.log("buffer: ", buf)

  var config = {
    method: 'post',
    url: 'http://172.18.0.2:1317/airequest/aireq',
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
});
