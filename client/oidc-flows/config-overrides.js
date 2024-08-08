const path = require('path');

module.exports = {
  webpack: function override(config, env) {
    if (env === 'production') {
      config.output.filename = 'static/js/main.js';
    }
    return config;
  },
  paths: function (paths, env) {
    paths.appBuild = path.resolve(__dirname, 'oidc-flows');
    return paths;
  }
};