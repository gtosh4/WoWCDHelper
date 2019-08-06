module.exports = {
  devServer: {
    proxy: {
      '/': {
        target: 'http://localhost:8113',
        ws: false,
      }
    }
  },
  configureWebpack: {
    devtool: 'source-map'
  },
}
