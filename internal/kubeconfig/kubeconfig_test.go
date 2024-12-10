package kubeconfig

import (
	"encoding/json"
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
