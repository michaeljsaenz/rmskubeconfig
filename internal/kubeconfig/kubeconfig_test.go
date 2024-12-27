package kubeconfig

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/michaeljsaenz/rmskubeconfig/internal/types"
	yaml "gopkg.in/yaml.v3"
)

func TestGetClusters_Success(t *testing.T) {
	// mock response data
	expectedClusters := []types.RMSCluster{
		{ID: "1", Name: "Cluster-1"},
		{ID: "2", Name: "Cluster-2"},
	}
	mockClusterResponse := types.RMSClusterResponse{Data: expectedClusters}

	// mock rms-api server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != clusterListPath {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockClusterResponse)
	}))
	defer mockServer.Close()

	clusters, err := GetClusters(mockServer.URL, "mockApiToken")
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	if !reflect.DeepEqual(clusters, expectedClusters) {
		t.Errorf("Expected %v, but got %v", expectedClusters, clusters)
	}

}
func TestGetClusters_Unauthorized(t *testing.T) {
	// mock rms-api server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer mockServer.Close()

	_, err := GetClusters(mockServer.URL, "mockApiToken")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	if !strings.Contains(err.Error(), "401") {
		t.Errorf("Expected 401 error, but got: %v", err)
	}

}

func TestGetClusters_DoRequestErrorNoHost(t *testing.T) {
	// invalid host (i.e., no host in URL)
	_, err := GetClusters("http://", "mockApiToken")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	var reqErr *types.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected custom RequestError, but got: %T", err)
	}

	if reqErr.Code != types.ErrRequestCode {
		t.Errorf("expected error code %d, but got: %d", types.ErrRequestCode, reqErr.Code)
	}

}

func TestGetClusters_NewRequestInvalidScheme(t *testing.T) {
	// missing protocol scheme (i.e., missing http/https)
	_, err := GetClusters("://missing-scheme", "mockApiToken")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	var reqErr *types.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected custom RequestError, but got: %T", err)
	}

	if reqErr.Code != types.ErrRequestCode {
		t.Errorf("expected error code %d, but got: %d", types.ErrRequestCode, reqErr.Code)
	}

}

func TestGetClusters_StatusNotFound(t *testing.T) {
	// mock rms-api server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer mockServer.Close()

	_, err := GetClusters(mockServer.URL, "mockApiToken")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expected 404 error, but got: %v", err)
	}

}

func TestGetClusters_ErrorDecodingResponse(t *testing.T) {
	// mock rms-api server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return non-JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer mockServer.Close()

	_, err := GetClusters(mockServer.URL, "mockApiToken")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	var reqErr *types.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected custom RequestError, but got: %T", err)
	}

	if reqErr.Code != types.ErrRequestCode {
		t.Errorf("Expected error code %d, but got: %d", types.ErrRequestCode, reqErr.Code)
	}

}

func TestGenerateCombinedKubeconfig_Success(t *testing.T) {
	// Mock data for the kubeconfig response
	mockKubeconfigResponseCluster1 := types.KubeconfigResponse{
		Config: `
clusters:
- name: cluster1
  cluster:
    server: https://cluster1.test
users:
- name: user1
  user:
    token: token1
contexts:
- name: context1
  context:
    cluster: cluster1
    user: user1`,
	}
	mockKubeconfigResponseCluster2 := types.KubeconfigResponse{
		Config: `
clusters:
- name: cluster2
  cluster:
    server: https://cluster2.test
users:
- name: user2
  user:
    token: token2
contexts:
- name: context2
  context:
    cluster: cluster2
    user: user2`,
	}

	// Mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Expected POST request, got %s", r.Method)
		}

		// extract the clusterID from URL
		urlPath := r.URL.Path
		clusterID := strings.TrimPrefix(urlPath, clusterListPath)

		// Simulate kubeconfig response
		switch clusterID {
		case "cluster1":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockKubeconfigResponseCluster1)
		case "cluster2":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockKubeconfigResponseCluster2)
		default:
			http.Error(w, "cluster not found", http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// create temp dir
	tempDir, err := os.MkdirTemp("", "rmskubeconfig-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = GenerateCombinedKubeconfig(mockServer.URL, "mock-token", tempDir, []string{"cluster1", "cluster2"})
	if err != nil {
		t.Fatalf("Function returned an error: %v", err)
	}

	// validate the returned output (combined kubeconfig)
	outputPath := tempDir + "/config"
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var combinedKubeconfig types.Kubeconfig
	err = yaml.Unmarshal(output, &combinedKubeconfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal combined kubeconfig: %v", err)
	}

	// assertions
	if len(combinedKubeconfig.Clusters) != 2 {
		t.Errorf("Expected 2 clusters, got %d", len(combinedKubeconfig.Clusters))
	}
	if len(combinedKubeconfig.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(combinedKubeconfig.Users))
	}
	if len(combinedKubeconfig.Contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(combinedKubeconfig.Contexts))
	}
}

func TestGenerateCombinedKubeconfig_ClusterNotFound(t *testing.T) {
	// Mock HTTP server - kubeconfig response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	err := GenerateCombinedKubeconfig(mockServer.URL, "mock-token", "", []string{"cluster-does-not-exist"})
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	var reqErr *types.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected custom RequestError, but got: %T", err)
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expected 404 error, but got: %v", err)
	}
}

func TestGenerateCombinedKubeconfig_NewRequestInvalidScheme(t *testing.T) {
	// missing protocol scheme (i.e., missing http/https)
	err := GenerateCombinedKubeconfig("://missing-scheme", "mock-token", "", []string{"cluster-does-not-exist"})

	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	var reqErr *types.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected custom RequestError, but got: %T", err)
	}

	if reqErr.Code != types.ErrRequestCode {
		t.Errorf("expected error code %d, but got: %d", types.ErrRequestCode, reqErr.Code)
	}

}

func TestGenerateCombinedKubeconfig_DoRequestErrorNoHost(t *testing.T) {
	// invalid host (i.e., no host in URL)
	err := GenerateCombinedKubeconfig("https://", "mock-token", "", []string{"cluster-does-not-exist"})
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	var reqErr *types.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected custom RequestError, but got: %T", err)
	}

	if reqErr.Code != types.ErrRequestCode {
		t.Errorf("expected error code %d, but got: %d", types.ErrRequestCode, reqErr.Code)
	}
}
