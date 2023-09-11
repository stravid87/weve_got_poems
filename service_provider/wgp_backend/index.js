var express = require('express');
var path = require('path');

// INIT
const PORT = 8080;
const app = express();

// MIDDLEWARE
app.use(express.static("../wgp_frontend/dist"));

// ROUTES
app.get('/', (req, res) => {
    res.sendFile('index.html')
  })

// lISTEN
app.listen(PORT, () => {
    console.log(`Example app listening on port ${PORT}`)
})