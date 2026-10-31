package gotools

import (
	"debug/buildinfo"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

// TestHelmComesFromTheToolsInput is the proof of go-check's tools input:
// helm is on PATH in this job's own $RUNNER_TEMP/bin, built from the module
// version mise.toml pins. A helm that the runner's image or an earlier job
// left anywhere else on PATH fails it.
func TestHelmComesFromTheToolsInput(t *testing.T) {
	temp := os.Getenv("RUNNER_TEMP")
	if temp == "" {
		t.Skip("RUNNER_TEMP is unset: this test runs in go-check's test job")
	}
	manifest, err := os.ReadFile("mise.toml")
	if err != nil {
		t.Fatal(err)
	}
	pin := regexp.MustCompile(`(?m)^helm = "([^"]+)"$`).FindSubmatch(manifest)
	if pin == nil {
		t.Fatal("mise.toml pins no helm")
	}
	path, err := exec.LookPath("helm")
	if err != nil {
		t.Fatalf("helm is not on PATH: %v", err)
	}
	if want := filepath.Join(temp, "bin", "helm"); path != want {
		t.Fatalf("helm resolves to %s, want %s, where the tools input installs it", path, want)
	}
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		t.Fatalf("%s carries no Go build info: %v", path, err)
	}
	if want := "v" + string(pin[1]); info.Main.Path != "helm.sh/helm/v4" || info.Main.Version != want {
		t.Fatalf("%s is %s %s, want helm.sh/helm/v4 %s from mise.toml", path, info.Main.Path, info.Main.Version, want)
	}
}
