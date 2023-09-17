const fs = require('fs');
const crypto = require("crypto").webcrypto;
globalThis.crypto = crypto;
require('./wasm_exec.js');


const wasmModule = fs.readFileSync('../../layer8/middleware/middleware.wasm');
const go = new Go();
const importObject = go.importObject;
WebAssembly.instantiate(wasmModule, importObject).then((results) => {
    const instance = results.instance
    go.run(instance);
    console.log("1 + 2 = ", addTwoNumbers(1,2))
});

module.exports = function Layer8WASM(req, res, next) {
    reqProxyLogger(req, res, next);
};


