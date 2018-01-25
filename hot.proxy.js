var webpack = require('webpack');
var webpackDevMiddleware = require('webpack-dev-middleware');
var webpackHotMiddleware = require('webpack-hot-middleware');
var proxy = require('proxy-middleware');
var config = require('./webpack.config');

var port = +(process.env.PORT || 5000);

hotBundle = ['webpack-hot-middleware/client?http://localhost:' + port].concat(config.entry.bundle);

config.entry = {
    bundle: hotBundle
};

config.plugins.push(
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin()
);

config.devtool = 'eval-source-map';

var app = new require('express')();

var compiler = webpack(config);
app.use(webpackDevMiddleware(compiler, {
    noInfo: true,
    publicPath: config.output.publicPath
}));
app.use(webpackHotMiddleware(compiler));
app.use(proxy('http://localhost:' + port));

port++;

app.listen(port, function(error) {
    if (error) {
        console.error(error);
    } else {
        console.info('==> 🌎  Listening on port %s. Open up http://localhost:%s/ in your browser.', port, port);
    }
});
