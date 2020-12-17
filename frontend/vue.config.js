module.exports = {
  devServer: {
    proxy: {
      '/': {
        target: 'http://localhost:8113',
      }
    }
  },
  configureWebpack: {
    devtool: 'source-map'
  },
}
