const express = require("express");
const controller = require("./kycController.js");
const { check } = require("express-validator/check");
const multer = require('multer');
const fs = require('fs');

// destination and name of the image file 
const storage = multer.diskStorage({
  destination: function (req, file, cb) {
    // if the directory does not exist then we create a new one
    fs.mkdir('./.uploads/',(err)=>{
      cb(null, './.uploads');
   });
  },
  filename: function (req, file, cb) {
    const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9)
    cb(null, file.originalname)
  }
})
const upload = multer({ storage: storage })

const router = express.Router();

router.post("/img_kyc", upload.single('image'), [
  check("oscript_name").isLength({ min: 3, max: 99 }),
  check("input").isBase64(),
  check("expected_output").isBase64(),
  check("fees").isNumeric(),
  check("validator_count").isNumeric().isInt({ min: 1 })
], controller.getKYC);

module.exports = router;