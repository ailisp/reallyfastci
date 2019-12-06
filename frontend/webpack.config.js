const path = require('path');
const webpack = require('webpack');
const api_url = process.env.NODE_ENV == 'production' ? process.env.REALLYFASTCI_BACKEND : 'http://localhost:1323'
module.exports = {
  entry: {
    'dist/app': './src/main.ts',
  },
  output: {
    filename: '[name].js',
    path: path.resolve(__dirname)
  },
  plugins: [
    new webpack.DefinePlugin({
      "webpackenv.API_URL": JSON.stringify(api_url)
    })
  ],
  resolve: {
    extensions: ['.ts', '.tsx', '.js']
  },
  module: {
    rules: [
      { test: /.tsx?$/, loader: 'ts-loader' },
      { test: /\.js$/, use: ["source-map-loader"], enforce: "pre" }
    ]
  },
  devServer: {
    open: true
  },
  devtool: 'source-map'
}