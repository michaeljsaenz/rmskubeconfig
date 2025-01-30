package rmskubeconfig

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/michaeljsaenz/rmskubeconfig/internal/kubeconfig"
	"github.com/michaeljsaenz/rmskubeconfig/internal/types"
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
		t.Errorf("expected outputPath %q, got %q", expectedTempDir, c.outputPath)
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
func TestRMSUrl(t *testing.T) {

	expectedUrl := "https//test.url"
	config := &Config{
		rmsUrl: expectedUrl,
	}

	actualUrl := config.RMSUrl()

	if expectedUrl != actualUrl {
		t.Errorf("RMSUrl() expected %q; got %q", expectedUrl, actualUrl)
	}

}

func TestApiToken(t *testing.T) {

	expectedApiToken := "token-test:1234"
	config := &Config{
		apiToken: expectedApiToken,
	}

	actualApiToken := config.ApiToken()

	if expectedApiToken != actualApiToken {
		t.Errorf("ApiToken() expected %q; got %q", expectedApiToken, actualApiToken)
	}

}
func TestOutputPath(t *testing.T) {

	expectedOutputPath := "/test/path/"
	config := &Config{
		outputPath: expectedOutputPath,
	}

	actualOutputPath := config.OutputPath()

	if expectedOutputPath != actualOutputPath {
		t.Errorf("OutputPath() expected %q; got %q", expectedOutputPath, actualOutputPath)
	}

}

func TestRun_Success(t *testing.T) {
	// mock response data
	expectedClusters := []types.RMSCluster{
		{ID: "1", Name: "Cluster-1"},
	}
	mockClusterResponse := types.RMSClusterResponse{Data: expectedClusters}

	// mock data for the kubeconfig response
	mockKubeconfigResponseCluster := types.KubeconfigResponse{
		Config: `
clusters:
- name: cluster1
  cluster:
    server: https://cluster1.test
users:
- name: user1
  user:
    token: token
contexts:
- name: context1
  context:
    cluster: cluster1
    user: user1`,
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if action := r.URL.Query().Get("action"); action == kubeconfig.GenerateKubeconfigUrlAction {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockKubeconfigResponseCluster)
		} else if r.URL.Path == kubeconfig.ClusterListPath {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockClusterResponse)
		}
	}))
	defer mockServer.Close()

	c := &Config{
		rmsUrl:     mockServer.URL,
		apiToken:   "token-test:test",
		outputPath: t.TempDir(),
	}

	c.Run()

	_, err := os.ReadFile(c.outputPath + "/config")
	if err != nil {
		t.Errorf("failed to read combined kubeconfig config file, error: %v", err)
	}
}

func TestRun_DefaultOutputPath(t *testing.T) {
	// mock response data
	expectedClusters := []types.RMSCluster{
		{ID: "1", Name: "Cluster-1"},
	}
	mockClusterResponse := types.RMSClusterResponse{Data: expectedClusters}

	// mock data for the kubeconfig response
	mockKubeconfigResponseCluster := types.KubeconfigResponse{
		Config: `
clusters:
- name: cluster1
  cluster:
    server: https://cluster1.test
users:
- name: user1
  user:
    token: token
contexts:
- name: context1
  context:
    cluster: cluster1
    user: user1`,
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if action := r.URL.Query().Get("action"); action == kubeconfig.GenerateKubeconfigUrlAction {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockKubeconfigResponseCluster)
		} else if r.URL.Path == kubeconfig.ClusterListPath {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockClusterResponse)
		}
	}))
	defer mockServer.Close()

	c := &Config{
		rmsUrl:   mockServer.URL,
		apiToken: "token-test:test",
	}

	c.Run()

	_, err := os.ReadFile(c.outputPath + "/config")
	if err != nil {
		t.Errorf("failed to read combined kubeconfig config file, error: %v", err)
	}
}

func TestRun_CombinedKubconfigError(t *testing.T) {
	// mock response data
	expectedClusters := []types.RMSCluster{
		{ID: "1", Name: "Cluster-1"},
	}
	mockClusterResponse := types.RMSClusterResponse{Data: expectedClusters}

	// mock data for the kubeconfig response
	mockKubeconfigResponseCluster := types.KubeconfigResponse{
		Config: `
clusters:
- name: cluster1
  cluster:
    server: https://cluster1.test
users:
- name: user1
  user:
    token: token
contexts:
- name: context1
  context:
    cluster: cluster1
    user: user1`,
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if action := r.URL.Query().Get("action"); action == kubeconfig.GenerateKubeconfigUrlAction {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockKubeconfigResponseCluster)
		} else if r.URL.Path == kubeconfig.ClusterListPath {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockClusterResponse)
		}
	}))
	defer mockServer.Close()

	c := &Config{
		rmsUrl:     mockServer.URL,
		apiToken:   "token-test:test",
		outputPath: t.TempDir() + "/invalid/path/",
	}

	expectedError := "no such file or directory"

	err := c.Run()
	if err != nil && !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error message to contain %q, but got: %v", expectedError, err)
	}

	if err == nil {
		t.Fatalf("expected error, but got: %v", err)
	}
}
