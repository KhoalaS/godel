package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wasmerio/wasmer-go/wasmer"
)

type WasmRunner struct {
}

func (wr *WasmRunner) Run(input string) (any, error) {
	outputWasmFile, err := wr.GenerateWasm(input, "./")
	if err != nil {
		return nil, err
	}

	return wr.ExecuteWasm(outputWasmFile)
}

func (wr *WasmRunner) ExecuteWasm(wasmFilePath string) (any, error) {
	wasmBytes, err := os.ReadFile(wasmFilePath)
	if err != nil {
		return nil, err
	}

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	module, _ := wasmer.NewModule(store, wasmBytes)

	// Instantiates the module
	importObject := wasmer.NewImportObject()
	instance, _ := wasmer.NewInstance(module, importObject)

	// Gets the `sum` exported function from the WebAssembly instance.
	hello, err := instance.Exports.GetFunction("hello")
	if err != nil {
		return nil, err
	}

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	result, err := hello()
	if err != nil {
		return nil, err
	}

	fmt.Println(result)
	return result, nil
}

func (wr *WasmRunner) GenerateWasm(input string, outputDir string) (string, error) {

	var _outPutDir string
	if outputDir == "" {
		_outPutDir = os.TempDir()
	} else {
		_outPutDir = outputDir
	}

	jsFile, err := os.CreateTemp(_outPutDir, "input-*.js")
	if err != nil {
		return "", err
	}

	_, err = jsFile.WriteString(input)
	if err != nil {
		return "", err
	}

	outputWasmFile := filepath.Join(_outPutDir, "out.wasm")

	porfCommand := exec.Command("porf", "wasm", jsFile.Name(), outputWasmFile)
	err = porfCommand.Run()
	if err != nil {
		return "", err
	}

	return outputWasmFile, nil
}
