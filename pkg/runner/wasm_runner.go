package runner

import (
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/KhoalaS/godel/pkg/custom_error"
	"github.com/rs/zerolog/log"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type WasmRunner struct {
}

func (wr *WasmRunner) Run(input string) (any, error) {
	outputWasmFile, err := wr.GenerateWasm(input, "")
	if err != nil {
		return nil, custom_error.FromError(err, WASM_GENERATION_ERROR_CODE, "runner")
	}

	return wr.ExecuteWasm(outputWasmFile)
}

func (wr *WasmRunner) ExecuteWasm(wasmFilePath string) (any, error) {
	defer os.Remove(wasmFilePath)

	wasmBytes, err := os.ReadFile(wasmFilePath)
	if err != nil {
		return nil, custom_error.FromError(err, FILE_READ_ERROR_CODE, "runner")
	}

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	module, _ := wasmer.NewModule(store, wasmBytes)

	importObject := wasmer.NewImportObject()
	instance, _ := wasmer.NewInstance(module, importObject)

	memory, _ := instance.Exports.GetMemory("$")
	mainFn, _ := instance.Exports.GetFunction("main")

	results, err := mainFn(0.0, 0, 0.0, 0)
	if err != nil {
		log.Err(err).Msg("Failed to call main")
		return nil, custom_error.FromError(err, WASM_EXECUTION_ERROR_CODE, "runner")
	}

	res := results.([]any)
	if len(res) != 2 {
		log.Fatal().Int("count", len(res)).Msg("Unexpected return value count")
	}

	offset := int32(res[0].(float64))
	tag := res[1].(int32)

	if tag == 1 {
		return offset, nil
	}

	log.Debug().Int32("offfset", offset).Int32("tag", tag).Msg("Raw return values")

	data := memory.Data()

	length := int32(binary.LittleEndian.Uint32(data[offset : offset+4]))

	strStart := offset + 4
	strBytes := data[strStart : strStart+length]

	log.Debug().Str("result", string(strBytes)).Send()

	return string(strBytes), nil
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
		return "", custom_error.FromError(err, FILE_CREATE_ERROR_CODE, "runner")
	}

	_, err = jsFile.WriteString(input)
	if err != nil {
		return "", custom_error.FromError(err, FILE_WRITE_ERROR_CODE, "runner")
	}

	jsFile.Close()
	defer os.Remove(jsFile.Name())

	outputWasmFile := filepath.Join(_outPutDir, "out.wasm")

	porfCommand := exec.Command("porf", "wasm", "--module", jsFile.Name(), outputWasmFile)
	err = porfCommand.Run()
	if err != nil {
		return "", custom_error.FromError(err, PORF_ERROR_CODE, "runner")
	}

	return outputWasmFile, nil
}
