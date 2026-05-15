package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestMain_MissingConfig verifies the binary exits non-zero when the config
// file cannot be found. This is an integration-style smoke test.
func TestMain_MissingConfig(t *testing.T) {
	if os.Getenv("PULSECTL_RUN_INTEGRATION") != "1" {
		t.Skip("set PULSECTL_RUN_INTEGRATION=1 to run")
	}

	tmp := t.TempDir()
	binary := filepath.Join(tmp, "pulsectl")

	build := exec.Command("go", "build", "-o", binary, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd := exec.Command(binary, "-config", filepath.Join(tmp, "nonexistent.yaml"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected non-zero exit for missing config, got nil")
		}
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("process did not exit in time")
	}
}
