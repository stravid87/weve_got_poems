function loadWebAssembly() {
    const wasmModule = fs.readFileSync('../../layer8/middleware/middleware.wasm');
      const go = new Go();
      const importObject = go.importObject;
      WebAssembly.instantiate(wasmModule, importObject).then((results) => {
        const instance = results.instance
        go.run(instance);
        console.log("loadWebAssembly run")
      });
  }

