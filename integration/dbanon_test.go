package main

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

var binaryName = "dbanon"

var version string

func TestDbanon(t *testing.T) {
	os.Chdir("..")
	dir, _ := os.Getwd()
	
	cmd1 := exec.Command(path.Join(dir, binaryName))
	out1, _ := cmd1.CombinedOutput()
	res1 := strings.TrimSpace(string(out1))
	if res1 != "You must specify a config" {
		t.Errorf("Got %s expected missing config error", res1)
	}

	versCmd := exec.Command("git", "describe", "--tags")
	versOut, _ := versCmd.CombinedOutput()
	versRes := strings.TrimSpace(string(versOut))

	cmd2 := exec.Command(path.Join(dir, binaryName), "-version")
	out2, _ := cmd2.CombinedOutput()
	res2 := strings.TrimSpace(string(out2))
	if res2 != versRes {
		t.Errorf("Got %s expected %s", res2, versRes)
	}


	cmdStr3 := "cat integration/magento_raw.sql | ./dbanon -config=magento2"
	cmd3 := exec.Command("bash", "-c", cmdStr3)
	out3, _ := cmd3.CombinedOutput()
	res3 := strings.TrimSpace(string(out3))
	if strings.Contains(res3, "adminFirstName") {
		t.Error("Expected no adminFirstName")
	}
	if strings.Contains(res3, "customer1FirstName") {
		t.Error("Expected no customer1FirstName")
	}
	if strings.Contains(res3, "customer2LastName") {
		t.Error("Expected no customer2LastName")
	}
}