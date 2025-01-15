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
	}

	mainJSContents := string(fileBytes)

	ctx := v8.NewContext()
	ctx.RunScript(mainJSContents, "main.js")
	val, err := ctx.RunScript("out", "test-get-value.js") // any functions previously added to the context can be called
	if err != nil {
		panic(err)
	}
	fmt.Println(val.String())
}
