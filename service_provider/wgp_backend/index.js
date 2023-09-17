var express = require('express');
var cors = require('cors');
const Layer8WASM = require("../../layer8/middleware/middleware.js")

// INIT
const PORT = 3000;
const app = express();


app.use(express.static("../wgp_frontend/dist"));
app.use(cors())
app.use(express.json())
app.use(Layer8WASM)

// ROUTES


// TOMORROWS LABOUR: GETTING / WORKS BUT THE API ENDPOINT NO?

app.get('/', (req, res) => {
    console.log("Arrived in '/home'")
    res.sendFile('index.html')
  })

app.get('api/v1/ping-service-provider', (req, res) => {
  console.log("req.body", req.body)
  res.send("Response from '/ping-service-provider'")
})


// lISTEN
app.listen(PORT, () => {
    console.log(`Example app listening on port ${PORT}`)
})
  


