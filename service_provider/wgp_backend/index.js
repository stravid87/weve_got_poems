var express = require('express');
var path = require('path');

// INIT
const PORT = 8080;
const app = express();

// MIDDLEWARE
app.use(express.static("../wgp_frontend/dist"));

// ROUTES
app.get('/', (req, res) => {
    console.log("Arrived in '/home'")
    res.sendFile('index.html')
  })

app.get('/route2', (req, res) => {
  console.log("Arrived in '/route2'")
  res.send("Response from '/route2'")
})

// lISTEN
app.listen(PORT, () => {
    console.log(`Example app listening on port ${PORT}`)
})