package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	v8 "rogchap.com/v8go"
)

const port = 8080

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/run", handleRun)

	log.Printf("Listening on port, %d", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type RunRequest struct {
	FunctionDefinition string
	Argument           string
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	log.Print("/run")
	var params RunRequest
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := runCode(params.FunctionDefinition, params.Argument)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(output))
}

func runCode(jsFunction string, arg string) (string, error) {
	fileBytes, err := os.ReadFile("dist/main.js")
	if err != nil {
		return "", err
	}
	mainJSContents := string(fileBytes)

	ctx := v8.NewContext()
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	// Run the bundled javascript file, which includes npm dependencies and defines a global variable "out"
	_, err = ctx.RunScript(mainJSContents, "main.js")
	if err != nil {
		return "", err
	}

	_, err = ctx.RunScript(jsFunction, "define-user-function.js")
	if err != nil {
		return "", err
	}

	execute, err := ctx.Global().Get("execute")
	if err != nil {
		return "", err
	}
	executeFn, err := execute.AsFunction()
	if err != nil {
		return "", fmt.Errorf("input must include a function called 'execute'")
	}

	argCtx := v8.NewContext()
	argIso := argCtx.Isolate()
	defer argIso.Dispose()
	defer argCtx.Close()

	_, err = argCtx.RunScript(fmt.Sprintf("arg = %s", arg), "get-arg-value.js")
	if err != nil {
		return "", err
	}
	argVal, err := argCtx.Global().Get("arg")
	if err != nil {
		return "", err
	}

	outputVal, err := executeFn.Call(v8.Undefined(iso), argVal)
	if err != nil {
		return "", err
	}

	if !outputVal.IsObject() {
		return outputVal.String(), nil

	}

	outputJSON, err := v8.JSONStringify(ctx, outputVal)
	if err != nil {
		return "", err
	}
	return outputJSON, nil
}
