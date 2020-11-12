var axios = require('axios');
var FormData = require('form-data');
var fs = require('fs');
var data = new FormData();
data.append('file', fs.createReadStream('../images/photo_2020-11-12_11-38-52.jpg'));

var config = {
  method: 'post',
  url: 'http://164.90.180.95:5001/api/v0/add',
  headers: { 
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