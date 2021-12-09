package main

import (
	"fmt"
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

	pwd, _ := os.Getwd()
	configPath := pwd + "/testdata/magento2.yml"

	fmt.Printf("cat integration/magento_raw.sql | ./dbanon -config=%s", configPath)

	cmdStr3 := fmt.Sprintf("cat integration/magento_raw.sql | ./dbanon -config=%s", configPath)
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

	cmdStr4 := fmt.Sprintf("cat integration/magento_eav_before.sql | ./dbanon -config=%s map-eav", configPath)
	cmd4 := exec.Command("bash", "-c", cmdStr4)
	out4, _ := cmd4.CombinedOutput()
	res4 := strings.TrimSpace(string(out4))
	if !strings.Contains(res4, "\"5\": firstname") {
		t.Errorf("First name not processed correctly")
	}

	if !strings.Contains(res4, "\"28\": street") {
		t.Errorf("Street not processed correctly")
	}

	cmdStr5 := "cat integration/laravel.sql | ./dbanon -config=integration/laravel.yml"
	cmd5 := exec.Command("bash", "-c", cmdStr5)
	out5, _ := cmd5.CombinedOutput()
	res5 := strings.TrimSpace(string(out5))
	if strings.Contains(res5, "Bob Smith") {
		t.Error("Expected no Bob Smith")
	}
}
