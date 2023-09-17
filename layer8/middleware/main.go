package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("addTwoNumbers", js.FuncOf(addTwoNumbers))
	js.Global().Set("reqProxyLogger", js.FuncOf(reqProxyLogger))
	<-c
}

func general_Form_WASM_Promise() interface{} {
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		//reject := args[1]
		go func() {
			// Main funciont body
			resolve.Invoke()
			//reject.Invoke()
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func reqProxyLogger(this js.Value, args []js.Value) interface{} {
	//request := args[0]
	//response := args[1]
	next := args[2]

	fmt.Println("Request has transitted the middleware.")

	next.Invoke()

	fmt.Println("up and down")

	return nil
}

func addTwoNumbers(this js.Value, args []js.Value) interface{} {
	a := args[0].Int()
	b := args[1].Int()
	sum := a + b
	return js.ValueOf(sum)
}
