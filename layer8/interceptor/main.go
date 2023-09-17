package main

import (
	"fmt"
	"io"
	"net/http"
	"syscall/js"
)

func main() {
	fmt.Println("WASM interceptor loaded and locked")
	c := make(chan struct{}, 0)
	js.Global().Set("addTwoNumbers", js.FuncOf(addTwoNumbers))
	js.Global().Set("reachOutToBackend", js.FuncOf(reachOutToBackend))
	js.Global().Set("pingSlave", js.FuncOf(pingSlave))
	js.Global().Set("pingServiceProvider", js.FuncOf(pingServiceProvider))
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

func addTwoNumbers(this js.Value, args []js.Value) interface{} {
	a := args[0].Int()
	b := args[1].Int()
	sum := a + b
	return js.ValueOf(sum)
}

func pingSlave(this js.Value, args []js.Value) interface{} {
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			response, err := http.Get("http://localhost:8000/api/v1/ping")
			if err != nil {
				fmt.Println("Error calling 'pingSlave'", err.Error())
				reject.Invoke()
			}
			if response == nil || response.Body == nil {
				fmt.Println("Response of response body from pingSlave is nil")
			}

			defer response.Body.Close()

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading resp.Body:", err.Error())
				reject.Invoke(js.ValueOf(err.Error()))
				return
			}

			// Checking the HTTP status code
			if response.StatusCode != http.StatusOK {
				fmt.Println("Server returned non-OK status: ", response.Status)
				reject.Invoke(js.ValueOf(string(responseBody)))
				return
			}

			response_string := string(responseBody)
			if err != nil {
				fmt.Println("Error decoding resp.Body:", err.Error())
				reject.Invoke(js.ValueOf(err.Error()))
			}

			fmt.Println("Inside WASM: ", response_string)
			resolve.Invoke(response_string)
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func pingServiceProvider(this js.Value, args []js.Value) interface{} {
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			response, err := http.Get("http://localhost:3000/api/v1/ping-service-provider")
			if err != nil {
				fmt.Println("Error calling 'pingServiceProvider'", err.Error())
				reject.Invoke()
			}
			if response == nil || response.Body == nil {
				fmt.Println("Response of response body from ping-service-provider is nil")
			}

			defer response.Body.Close()

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading resp.Body:", err.Error())
				reject.Invoke(js.ValueOf(err.Error()))
				return
			}

			// Checking the HTTP status code
			if response.StatusCode != http.StatusOK {
				fmt.Println("Server returned non-OK status: ", response.Status)
				reject.Invoke(js.ValueOf(string(responseBody)))
				return
			}

			response_string := string(responseBody)
			if err != nil {
				fmt.Println("Error decoding resp.Body:", err.Error())
				reject.Invoke(js.ValueOf(err.Error()))
			}

			fmt.Println("Inside WASM: ", response_string)
			resolve.Invoke(response_string)
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func reachOutToBackend(this js.Value, args []js.Value) interface{} {
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			// options := js.ValueOf(map[string]interface{}{
			// 	"method":  "GET",
			// 	"headers": js.ValueOf(map[string]interface{}{}),
			// })
			// method := options.Get("method").String()
			// if method == "" {
			// 	method = "GET"
			// }
			// headers := options.Get("headers")
			// body := options.Get("body").String()
			// // setting the body to an empty string if it's undefined
			// if body == "<undefined>" {
			// 	body = ""
			// }
			// headersMap := make(map[string]string)
			// js.Global().Get("Object").Call("keys", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// 	headersMap[args[0].String()] = args[1].String() // forEach, headersMap[key] = value
			// 	return nil
			// }))
			// msg := strings.NewReader("Hello from wasm")
			response, err := http.Get("http://localhost:8080/route2")
			if err != nil {
				fmt.Println("Error calling reachOutToBackend", err.Error())
				reject.Invoke()
			}
			if response == nil || response.Body == nil {
				fmt.Println("Response of response body from route2 is nil")
			}

			defer response.Body.Close()

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading resp.Body:", err.Error())
				reject.Invoke(js.ValueOf(err.Error()))
				return
			}

			// Checking the HTTP status code
			if response.StatusCode != http.StatusOK {
				fmt.Println("Server returned non-OK status: ", response.Status)
				reject.Invoke(js.ValueOf(string(responseBody)))
				return
			}

			response_string := string(responseBody)
			if err != nil {
				fmt.Println("Error decoding resp.Body:", err.Error())
				reject.Invoke(js.ValueOf(err.Error()))
			}

			fmt.Println("Inside WASM: ", response_string)
			resolve.Invoke(response_string)
			//reject.Invoke()
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}
