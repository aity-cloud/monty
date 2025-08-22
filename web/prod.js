const path = require('path');
const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();

app.use(express.static('./dist', { redirect: false }));

app.use('/monty-api', createProxyMiddleware({
  target:       'http://localhost:8888',
  pathRewrite:  { '^/monty-api': '' },
  secure:       false,
  changeOrigin: true,
}));

app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, '/dist/index.html'), { root: '/' });
});

// Start the server
app.listen(3000, () => {
  console.log('Server started on port 3000');
});
