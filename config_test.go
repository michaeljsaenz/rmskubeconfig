package rmskubeconfig

import (
	"testing"
)

func TestNewConfig(t *testing.T) {

	config := NewConfig()

	if config == nil {
		t.Fatalf("expected a non-nil Config instance")
	}
	if config.rmsUrl != "" {
		t.Fatalf("expected rmsUrl to be an empty string, got %q", config.rmsUrl)

	}
	if config.apiToken != "" {
		t.Fatalf("expected apiToken to be an empty string, but was not")

	}
	if config.outputPath != "" {
		t.Errorf("expected outputPath to be an empty string, got %q", config.outputPath)
	}
}
func TestSetRMSUrl_InvalidUrl(t *testing.T) {

	config := NewConfig()

	err := config.SetRMSUrl("ftp://invalid-url//http://")

	if err == nil {
		t.Fatalf("expected error, but got: %v", err)
	}

}
func TestSetRMSUrl_ValidUrl(t *testing.T) {

	c := NewConfig()
	expectedUrl := "https://local.test"

	err := c.SetRMSUrl(expectedUrl)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c.rmsUrl != expectedUrl {
		t.Errorf("expected rmsUrl to be %q, but got: %q", expectedUrl, c.rmsUrl)
	}

}

func TestSetApiToken_InvalidToken(t *testing.T) {

	config := NewConfig()

	err := config.SetApiToken("must-start-with-token-")

	if err == nil {
		t.Fatalf("expected error, but got: %v", err)
	}

}

func TestSetApiToken_EmptyInput(t *testing.T) {

	c := NewConfig()
	emptyToken := ""

	err := c.SetApiToken(emptyToken)
	if err == nil {
		t.Errorf("expected error for empty token, but got none")
	}
	if c.apiToken != "" {
		t.Errorf("expected API token to remain empty, but was not")
	}

}

func TestSetApiToken_ValidInput(t *testing.T) {

	c := NewConfig()
	expectedToken := "token-test:test"

	err := c.SetApiToken(expectedToken)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.apiToken != expectedToken {
		t.Errorf("expected API token was not set correctly")
	}

}
func TestSetOutputPath_Success(t *testing.T) {

	c := NewConfig()
	expectedTempDir := t.TempDir()

	err := c.SetOutputPath(expectedTempDir)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c.outputPath != expectedTempDir {
		t.Errorf("expected outputPath was not set correctly")
	}

}
func TestSetOutputPath_MissingDirectoryError(t *testing.T) {

	c := NewConfig()
	expectedTempDir := t.TempDir()

	err := c.SetOutputPath(expectedTempDir + "/missing/directory")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

}
