package main

import (
	"fmt"
	"syscall/js"
)

func addTwoNumbers(this js.Value, args []js.Value) interface{} {
	a := args[0].Int()
	b := args[1].Int()
	sum := a + b
	return js.ValueOf(sum)
}

func main() {
	fmt.Println("Adding Numbers from WASM interceptor")
	c := make(chan struct{}, 0)
	js.Global().Set("addTwoNumbers", js.FuncOf(addTwoNumbers))
	<-c
}
