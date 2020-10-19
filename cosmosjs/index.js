const express = require("express");
const winston = require("winston");
const expressWinston = require("express-winston");

const server = express();
const http = require("http").Server(server);

require("dotenv").config();

server.use((req, res, next) => {
  const origin = req.get('origin');

  // TODO Add origin validation
  res.header('Access-Control-Allow-Origin', origin);
  res.header('Access-Control-Allow-Credentials', true);
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization, Cache-Control, Pragma');

  // intercept OPTIONS method
  if (req.method === 'OPTIONS') {
    res.sendStatus(204);
  } else {
    next();
  }
});

server.use(express.urlencoded({ extended: false }));
server.use(express.json({ limit: "5mb" }));

server.use(expressWinston.logger({
  transports: [
    new winston.transports.Console()
  ],
  format: winston.format.combine(
    winston.format.json(),
    winston.format.timestamp(),
    winston.format.colorize()
  )
}));

require("./routes.js")(server);

server.use(expressWinston.errorLogger({
  transports: [
    new winston.transports.Console()
  ],
  format: winston.format.combine(
    winston.format.json(),
    winston.format.timestamp(),
    winston.format.colorize()
  )
}));

let port = process.env.PORT;
if (port == null || port === "") {
  port = 8080;
}

const listen = http.listen(port, () => {
  console.log("server is running on port", listen.address().port);
});
