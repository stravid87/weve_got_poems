package main

import (
	"fmt"
	"interceptor/utils"
	"syscall/js"
)

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

func addTwoNumbers(this js.Value, args []js.Value) interface{} {
	a := args[0].Int()
	b := args[1].Int()
	sum := a + b
	return js.ValueOf(sum)
}

func reachOutToBackend() interface{} {
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		//reject := args[1]
		go func() {
			options := js.ValueOf(map[string]interface{}{
				"method":  "GET",
				"headers": js.ValueOf(map[string]interface{}{}),
			})
			method := options.Get("method").String()
			headers := options.Get("headers")
			body := options.Get("body").String()
			// setting the body to an empty string if it's undefined
			if body == "<undefined>" {
				body = ""
			}
			headersMap := make(map[string]string)
			js.Global().Get("Object").Call("keys", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				headersMap[args[0].String()] = args[1].String()
				return nil
			}))
			js.Global().Get("Object").Call("keys", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				headersMap[args[0].String()] = args[1].String()
				return nil
			}))
			msg := []byte("Hello from wasm")
			request := utils.NewRequest("GET", "http://localhost:8080/route2", headers, msg)
			resolve.Invoke()
			//reject.Invoke()
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func main() {
	fmt.Println("Adding Numbers from WASM interceptor")
	c := make(chan struct{}, 0)
	js.Global().Set("addTwoNumbers", js.FuncOf(addTwoNumbers))
	<-c
}
