package main

import (
	"fmt"
	"os"

	v8 "rogchap.com/v8go"
)

func main() {
	fileBytes, err := os.ReadFile("dist/main.js")
	if err != nil {
		fmt.Print(err)
		return
	}
	mainJSContents := string(fileBytes)

	ctx := v8.NewContext()

	// Run the bundled javascript file, which includes npm dependencies and defines a global variable "out"
	_, err = ctx.RunScript(mainJSContents, "main.js")
	if err != nil {
		fmt.Print(err)
		return
	}
	// Get the value of the global variable "out"
	val, err := ctx.RunScript("out", "test-get-value.js")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(val.String())
}
