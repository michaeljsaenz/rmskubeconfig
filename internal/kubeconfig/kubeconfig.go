package kubeconfig

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/michaeljsaenz/rmskubeconfig/internal/types"

	yaml "gopkg.in/yaml.v3"
)

const ClusterListPath string = "/v3/clusters/"
const GenerateKubeconfigUrlAction string = "generateKubeconfig"

// GetClusters retrieves a list of all clusters from RMS
func GetClusters(baseUrl, apiToken string) ([]types.RMSCluster, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseUrl+ClusterListPath, nil)
	if err != nil {
		return nil, &types.RequestError{
			Code:    types.ErrRequestCode,
			Message: fmt.Sprintf("error creating cluster request: %v", err),
		}
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, &types.RequestError{
			Code:    types.ErrRequestCode,
			Message: fmt.Sprintf("error fetching clusters: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &types.RequestError{
			Code:    types.ErrRequestCode,
			Message: fmt.Sprintf("unexpected response status fetching clusters: %v", resp.Status),
		}
	}

	var clusterResp types.RMSClusterResponse
	if err := json.NewDecoder(resp.Body).Decode(&clusterResp); err != nil {
		return nil, &types.RequestError{
			Code:    types.ErrRequestCode,
			Message: fmt.Sprintf("error decoding cluster response: %v", err),
		}
	}

	return clusterResp.Data, nil

}

// GenerateCombinedKubeconfig combines all generated kubeconfig files into one kubeconfig (config) file
func GenerateCombinedKubeconfig(baseUrl, apiToken, outputPath string, clusterIDs []string) error {
	client := &http.Client{}
	combinedKubeconfig := &types.Kubeconfig{
		APIVersion: "v1",
		Kind:       "Config",
	}

	for _, clusterID := range clusterIDs {

		url := fmt.Sprintf("%s%s%s?action=%s", baseUrl, ClusterListPath, clusterID, GenerateKubeconfigUrlAction)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return &types.RequestError{
				Code:    types.ErrRequestCode,
				Message: fmt.Sprintf("error creating generate kubeconfig request: %v", err),
			}
		}

		req.Header.Set("Authorization", "Bearer "+apiToken)

		resp, err := client.Do(req)
		if err != nil {
			return &types.RequestError{
				Code:    types.ErrRequestCode,
				Message: fmt.Sprintf("error fetching kubeconfig generate for cluster: %s, error: %v", clusterID, err),
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return &types.RequestError{
				Code:    types.ErrRequestCode,
				Message: fmt.Sprintf("Unexpected response status generating kubeconfig for cluster: %s (%v)", clusterID, resp.Status),
			}
		}

		var kubeconfigResp types.KubeconfigResponse
		if err := json.NewDecoder(resp.Body).Decode(&kubeconfigResp); err != nil {
			return &types.RequestError{
				Code:    types.ErrRequestCode,
				Message: fmt.Sprintf("error decoding generate kubeconfig response for cluster: %s, error: %v", clusterID, err),
			}
		}

		var kubeconfig types.Kubeconfig
		err = yaml.Unmarshal([]byte(kubeconfigResp.Config), &kubeconfig)
		if err != nil {
			return &types.RequestError{
				Code:    types.ErrRequestCode,
				Message: fmt.Sprintf("error unmarshaling YAML (generate kubeconfig response) for cluster: %s, error: %v", clusterID, err),
			}
		}

		combinedKubeconfig.Clusters = append(combinedKubeconfig.Clusters, kubeconfig.Clusters...)
		combinedKubeconfig.Users = append(combinedKubeconfig.Users, kubeconfig.Users...)
		combinedKubeconfig.Contexts = append(combinedKubeconfig.Contexts, kubeconfig.Contexts...)
	}

	err := createConfigFile(combinedKubeconfig, outputPath)
	if err != nil {
		return err
	}

	return nil

}

func createConfigFile(combinedKubeconfig *types.Kubeconfig, outputPath string) error {
	combinedKubeconfigYaml, _ := yaml.Marshal(combinedKubeconfig)

	err := os.WriteFile(outputPath+"/config", combinedKubeconfigYaml, 0644)
	if err != nil {
		return fmt.Errorf("error creating combined kubeconfig config file, error: %v", err)
	}

	return nil
}
