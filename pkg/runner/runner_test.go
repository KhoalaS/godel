package runner

import "testing"

func TestWasmRunner(t *testing.T) {
	runner := WasmRunner{}

	result, err := runner.Run("export function hello(){console.log('hello world');return 0}")
	if err != nil {
		t.Error(err)
	}

	if result != 0 {
		t.Error("wrong result", result)
	}

}
