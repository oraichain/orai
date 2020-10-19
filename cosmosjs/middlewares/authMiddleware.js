const jwt = require("jsonwebtoken");

module.exports = (req, res, next) => {
  try {
    const token = req.headers.authorization.split(" ")[1];
    const decoded = jwt.verify(token, process.env.SECRET_KEY, {
      algorithms: ["HS512"],
    });
    req.userData = decoded;
    console.log("req user data: ", req.userData);
    next();
  } catch (error) {
    console.log("error", error);
    console.log("req", req.headers.authorization);
    res.status(401).json({
      message: "Auth failed",
    });
  }
};
