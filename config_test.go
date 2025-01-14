package rmskubeconfig

import (
	"os"
	"testing"
)

func TestConfig_NewConfig(t *testing.T) {

	config := NewConfig()

	if config == nil {
		t.Fatalf("expected a non-nil Config instance")
	}
	if config.rmsUrl != "" {
		t.Fatalf("expected rmsUrl to be an empty string, got %q", config.rmsUrl)

	}
	if config.apiToken != "" {
		t.Fatalf("expected apiToken to be an empty string, got %q", config.apiToken)

	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if config.outputPath != cwd {
		t.Errorf("expected outputPath to be %q, got %q", cwd, config.outputPath)
	}
}
