package runner

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/KhoalaS/godel/pkg/custom_error"
)

type PorfforRunner struct {
}

func (r *PorfforRunner) Run(input string) (any, error) {
	return r.ExecuteJsInPorf(input, "")
}

func (wr *PorfforRunner) ExecuteJsInPorf(jsInput string, outputDir string) (any, error) {
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

	_, err = jsFile.WriteString(jsInput)
	if err != nil {
		return "", custom_error.FromError(err, FILE_WRITE_ERROR_CODE, "runner")
	}
	jsFile.Close()
	defer os.Remove(jsFile.Name())

	porfCommand := exec.Command("porf", "--secure", jsFile.Name())

	output, err := porfCommand.Output()
	if err != nil {
		return nil, custom_error.FromError(err, PORF_ERROR_CODE, "runner")
	}

	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	clean := re.ReplaceAllString(string(output), "")

	return strings.TrimSpace(clean), nil
}
