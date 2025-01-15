package rmskubeconfig

import (
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
	if config.outputPath != "" {
		t.Errorf("expected outputPath to be an empty string, got %q", config.outputPath)
	}
}
func TestConfig_SetInvalidUrl(t *testing.T) {

	config := NewConfig()

	err := config.SetRMSUrl("ftp://invalid-url//http://")

	if err == nil {
		t.Fatalf("expected error, but got: %v", err)
	}

}

func TestConfig_SetInvalidApiToken(t *testing.T) {

	config := NewConfig()

	err := config.SetApiToken("must-start-with-token-")

	if err == nil {
		t.Fatalf("expected error, but got: %v", err)
	}

}
