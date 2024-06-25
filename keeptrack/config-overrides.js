const path = require("path");

module.exports = function override(config, env) {
  if (env === "production") {
    config.output.filename = "static/js/main.js";
  }

  return config;
};
