const express = require("express");
const controller = require("./kycController.js");
const { check } = require("express-validator/check");
var multer  = require('multer')

// destination and name of the image file 
var storage = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, './.uploads')
  },
  filename: function (req, file, cb) {
    const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9)
    cb(null, file.originalname)
  }
})
const upload = multer({ storage: storage })

const router = express.Router();

router.post("/kyc_price", upload.single('face_img'), [
  check("oscript_name").isLength({ min: 3, max: 99 }),
  check("input").isBase64(),
  check("expected_output").isBase64(),
  check("fees").isNumeric(),
  check("validator_count").isNumeric().isInt({ min: 1 })
], controller.getClassification);

module.exports = router;