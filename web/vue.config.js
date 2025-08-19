// https://github.com/webpack/webpack/issues/13572#issuecomment-923736472
const crypto = require('crypto');
const cryptoOrigCreateHash = crypto.createHash;

crypto.createHash = algorithm => cryptoOrigCreateHash(algorithm === 'md4' ? 'sha256' : algorithm);

const config = require('@rancher/shell/vue.config');
const webpack = require('webpack');

const isStandalone = process.env.IS_STANDALONE === 'true';
let montyApi = process.env.MONTY_API || 'http://localhost:8888';

if (montyApi && !montyApi.startsWith('http')) {
  montyApi = `http://${ montyApi }`;
}

if (montyApi) {
  console.log(`MONTY API: ${ montyApi }`); // eslint-disable-line no-console
}

console.log(`IS STANDALONE`, isStandalone); // eslint-disable-line no-console

const baseConfig = config(__dirname, {
  excludes: [],
  proxies:  {
    '/monty-api': {
      secure:       false,
      target:       montyApi,
      pathRewrite:  { '^/monty-api': '' },
      ws:           true,
      changeOrigin: true,
    }
  }
  // excludes: ['fleet', 'example']
});

const baseConfigureWebpack = baseConfig.configureWebpack;

baseConfig.devServer.proxy = {
  '/monty-api': {
    secure:       false,
    target:       montyApi,
    pathRewrite:  { '^/monty-api': '' },
    ws:           true,
    changeOrigin: true,
  },
};

baseConfig.configureWebpack = (config) => {
  config.cache = { type: 'filesystem' };
  const comitHash = process.env.GITHUB_COMMIT_SHA || 'GITHUB_COMMIT_SHA not defined';

  config.plugins.push(new webpack.DefinePlugin({
    'process.env.isStandalone': JSON.stringify(isStandalone),
    'process.env.commitHash':   JSON.stringify(comitHash)
  }));

  baseConfigureWebpack(config);
};

// Makes the public path relative so that the <base> element will affect the assets.
if (!isStandalone) {
  baseConfig.publicPath = './';
}

// We need to add a custom script to the index in order to change how assets for the monty backendso we have to override the index.html
if (isStandalone) {
  baseConfig.pages.index.template = './pkg/monty/index.html';
}

baseConfig.productionSourceMap = false;

module.exports = baseConfig;
