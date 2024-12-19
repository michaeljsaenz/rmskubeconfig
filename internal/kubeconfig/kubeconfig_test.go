package kubeconfig

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/michaeljsaenz/rmskubeconfig/internal/types"
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

func TestGetClusters_InvalidUrlNoHost(t *testing.T) {
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

func TestGetClusters_InvalidScheme(t *testing.T) {
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
		w.Header().Set("Content-Type", "text/plain")
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
