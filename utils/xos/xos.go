package xos

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func ExecShell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func ExecBatch(s string, outputFile *string) (string, error) {
	var out bytes.Buffer

	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe", "/c", s)
	if outputFile != nil {
		stdout, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, os.FileMode.Perm(0600))
		if err != nil {
			return "", err
		}
		defer stdout.Close()
		cmd.Stdout = stdout
		cmd.Stderr = stdout
	} else {
		cmd.Stdout = &out
		cmd.Stderr = &out
	}

	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
