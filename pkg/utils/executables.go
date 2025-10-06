package utils

import "os/exec"

func ExecutableExists(name string) (bool, string) {
	absolutePath, err := exec.LookPath(name)
	if err != nil {
		return false, ""
	}

	return true, absolutePath
}
