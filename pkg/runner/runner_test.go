package runner

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestWasmRunner(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)

	var runner WasmRunner

	value := "'worlds'"
	inputJs := fmt.Sprintf("export function hello(s){return s};export function main(){return hello(%v)}", value)

	result, err := runner.Run(inputJs)
	if err != nil {
		t.Error(err)
	}

	parsedResult, ok := result.(string)
	if !ok {
		t.Error("result was not a string", result)
	}

	if parsedResult != "worlds" {
		t.Error("wrong result", result)
	}

	var intValue int32 = 100
	inputJs = fmt.Sprintf("export function hello(s){return s};export function main(){return hello(%d)}", intValue)


	result, err = runner.Run(inputJs)
	if err != nil {
		t.Error(err)
	}

	parsedIntResult, ok := result.(int32)
	if !ok {
		t.Error("result was not a string", result)
	}

	if parsedIntResult != intValue {
		t.Error("wrong result", result)
	}
}

func TestExecuteJsInPorf(t *testing.T) {
	jsInput := `function hello(s) {
    return s;
}
function main() {
    return hello({ name: "world" });
}

console.log(main());`

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)

	var runner PorfforRunner
	out, err := runner.ExecuteJsInPorf(jsInput, "./")
	if err != nil {
		t.Error(err)
	}

	val, ok := out.(string)
	if !ok {
		t.Error("output was not a string")
	}

	if val != "{\n  name: 'world'\n}" {
		t.Errorf("output is wrong got %s", val)
	}

}
