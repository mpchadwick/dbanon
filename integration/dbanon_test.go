package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

var binaryName = "dbanon"

func TestDbanon(t *testing.T) {
	dir, _ := os.Getwd()
	parent := filepath.Dir(dir)
	cmd := exec.Command(path.Join(parent, binaryName))
	output, _ := cmd.CombinedOutput()

	actual := strings.TrimSpace(string(output))

	if actual != "You must specify a config" {
		t.Errorf("Got %s expected missing config error", actual)
	}

}