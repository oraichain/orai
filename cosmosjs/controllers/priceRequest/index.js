const express = require("express");

const controller = require("./priceController.js");
const { check } = require("express-validator/check");

const router = express.Router();

router.post("/req_price", [
  check("oscript_name").isLength({ min: 3, max: 99 }),
  check("price").isBase64(),
  check("expected_price").isBase64(),
  check("fees").isNumeric()
], controller.getPrice);

module.exports = router;