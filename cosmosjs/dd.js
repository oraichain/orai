var axios = require('axios');
var FormData = require('form-data');
var fs = require('fs');
var data = new FormData();
data.append('image', fs.createReadStream('/home/duc/Pictures/photo_2020-11-12_11-38-52.jpg'));
data.append('oracle_script_name', 'oscript_classification');
data.append('fees', '45000orai');
data.append('from', 'orai18xc2v6uw7ph5q337awxmq43k5394gssvlx9q8e');
data.append('chain_id', 'Oraichain');
data.append('input', '\'\'');
data.append('expected_output', 'Samoyed');
data.append('validator_count', '1');

var config = {
  method: 'post',
  url: 'http://localhost:1317/airequest/aireq/kycreq',
  headers: { 
    'Content-Type': 'multipart/form-data', 
    ...data.getHeaders()
  },
  data : data
};

axios(config)
.then(function (response) {
  console.log(JSON.stringify(response.data));
})
.catch(function (error) {
  console.log(error);
});