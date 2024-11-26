package kubeconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/michaeljsaenz/rmskubeconfig/internal/types"

	yaml "gopkg.in/yaml.v3"
)

// GetClusters retrieves a list of all clusters from RMS
func GetClusters(baseURL, apiToken string) []types.RMSCluster {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL+"/v3/clusters", nil)
	if err != nil {
		log.Fatalf("Error creating cluster request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching clusters: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected response status fetching clusters: %d", resp.StatusCode)
	}

	var clusterResp types.RMSClusterResponse
	if err := json.NewDecoder(resp.Body).Decode(&clusterResp); err != nil {
		log.Fatalf("Error decoding cluster response: %v", err)
	}

	return clusterResp.Data

}

// GenerateCombinedKubeconfig combines all generated kubeconfig files into one kubeconfig (config) file
func GenerateCombinedKubeconfig(baseURL, apiToken, outputPath string, clusterIDs []string) error {

	client := &http.Client{}
	combinedKubeconfig := &types.Kubeconfig{
		APIVersion: "v1",
		Kind:       "Config",
	}

	for _, clusterID := range clusterIDs {

		url := fmt.Sprintf("%s/v3/clusters/%s?action=generateKubeconfig", baseURL, clusterID)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Fatalf("Error creating generate kubeconfig request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+apiToken)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making generate kubeconfig request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected response status generating kubeconfig: %d", resp.StatusCode)
		}

		var kubeconfigResp types.KubeconfigResponse
		if err := json.NewDecoder(resp.Body).Decode(&kubeconfigResp); err != nil {
			log.Fatalf("Error decoding kubeconfig response: %v", err)
		}

		var kubeconfig types.Kubeconfig
		err = yaml.Unmarshal([]byte(kubeconfigResp.Config), &kubeconfig)
		if err != nil {
			log.Fatalf("Error unmarshaling YAML (kubeconfig response): %v", err)
		}

		combinedKubeconfig.Clusters = append(combinedKubeconfig.Clusters, kubeconfig.Clusters...)
		combinedKubeconfig.Users = append(combinedKubeconfig.Users, kubeconfig.Users...)
		combinedKubeconfig.Contexts = append(combinedKubeconfig.Contexts, kubeconfig.Contexts...)
	}

	combinedKubeconfigYaml, err := yaml.Marshal(combinedKubeconfig)
	if err != nil {
		log.Fatalf("Failed to marshal combined kubeconfig YAML: %v", err)
	}

	createConfigFile(combinedKubeconfigYaml, outputPath)
	return nil

}

func createConfigFile(combinedKubeconfigYaml []byte, outputPath string) {
	err := os.WriteFile(outputPath+"/config", combinedKubeconfigYaml, 0644)
	log.Printf("Config file saved here: %s", outputPath+"/config")
	if err != nil {
		log.Fatalf("Error creating combined config file: %v", err)
	}
}
