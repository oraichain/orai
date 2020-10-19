module.exports = (server) => {
  server.use("/api/txs/", require("./controllers/priceRequest/index"));

  server.use("*", (req, res) => {
    console.log("Original url: ", req.originalUrl);
    res.status(404).json({ message: "Whoops, what are you looking for?" });
  });
};
