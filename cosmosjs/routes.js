module.exports = (server) => {
  server.use("/api/v1/txs/", require("./controllers/priceRequest/index"));

  server.use("/api/v1/txs/", require("./controllers/imgClassificationRequest/index"));

  server.use("/api/v1/txs/", require("./controllers/imgKYCRequest/index"));

  server.use("/api/v1/txs/", require("./controllers/imgOCRRequest/index"));

  server.use("*", (req, res) => {
    console.log("Original url: ", req.originalUrl);
    res.status(404).json({ message: "Whoops, what are you looking for?" });
  });
};
