package rmskubeconfig

import (
	"log"
	"os"

	"github.com/michaeljsaenz/rmskubeconfig/internal/kubeconfig"
)

// Config holds values for processing
type Config struct {
	RMSUrl     string
	APIToken   string
	OutputPath string
}

// NewConfig creates a new Config instance with default values
func NewConfig() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Config{
		RMSUrl:     "",
		APIToken:   "",
		OutputPath: cwd,
	}, nil
}

// SetRMSUrl sets RMS API URL
func (c *Config) SetRMSUrl(url string) {
	c.RMSUrl = url
}

// SetApiToken sets RMS API token
func (c *Config) SetApiToken(token string) {
	c.APIToken = token
}

// SetOutputPath sets path where to save config file
func (c *Config) SetOutputPath(path string) {
	c.OutputPath = path
}

// GetRMSUrl returns RMS API URL
func (c *Config) GetRMSUrl() string {
	return c.RMSUrl
}

// GetApiToken returns RMS API token
func (c *Config) GetApiToken() string {
	return c.APIToken
}

// GetOutputPath returns output file path
func (c *Config) GetOutputPath() string {
	return c.OutputPath
}

// Run executes the Config to generate combined kubeconfig (config) file
func (c *Config) Run() {
	clusters := kubeconfig.GetClusters(c.RMSUrl, c.APIToken)

	var clusterIDs []string
	for _, cluster := range clusters {
		clusterIDs = append(clusterIDs, cluster.ID)
	}

	err := kubeconfig.GenerateCombinedKubeconfig(c.RMSUrl, c.APIToken, clusterIDs)
	if err != nil {
		log.Fatalf("Error generating combined kubeconfig: %v", err)
	}
}
