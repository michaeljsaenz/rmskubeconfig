package rmskubeconfig

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/michaeljsaenz/rmskubeconfig/internal/kubeconfig"
)

// Config holds values for processing
type Config struct {
	rMSUrl     string
	aPIToken   string
	outputPath string
}

// NewConfig creates a new Config instance with default values
func NewConfig() *Config {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("NewConfig: failed to get current working directory: %v", err)
	}

	return &Config{
		rMSUrl:     "",
		aPIToken:   "",
		outputPath: cwd,
	}
}

// SetRMSUrl sets RMS API URL
func (c *Config) SetRMSUrl(url string) {
	// validate URL format
	regex := `^(https?:\/\/)?([\w\-]+(\.[\w\-]+)+)(:[0-9]{1,5})?(\/[^\s]*)?$`
	if match, _ := regexp.MatchString(regex, url); !match {
		log.Fatalf("SetRMSUrl: invalid RMS URL format: %s", url)
	}
	c.rMSUrl = url
}

// SetApiToken sets RMS API token
func (c *Config) SetApiToken(token string) {
	// validate token format
	regex := `^token-\w+:\w+`

	if match, _ := regexp.MatchString(regex, token); !match {
		log.Fatalf("SetApiToken: invalid API token format")
	}
	c.aPIToken = token
}

// SetOutputPath sets path where to save config file
func (c *Config) SetOutputPath(path string) {
	// validate path exists and ensure it's a directory
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		log.Fatalf("SetOutputPath: output path must be an existing directory: %s", path)
	}
	// convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("SetOutputPath: failed to resolve absolute path: %s", path)
	}
	c.outputPath = absPath
}

// GetRMSUrl returns RMS API URL
func (c *Config) GetRMSUrl() string {
	return c.rMSUrl
}

// GetApiToken returns RMS API token
func (c *Config) GetApiToken() string {
	return c.aPIToken
}

// GetOutputPath returns output file path
func (c *Config) GetOutputPath() string {
	return c.outputPath
}

// Run executes the Config to generate combined kubeconfig (config) file
func (c *Config) Run() {
	clusters := kubeconfig.GetClusters(c.rMSUrl, c.aPIToken)

	var clusterIDs []string
	for _, cluster := range clusters {
		clusterIDs = append(clusterIDs, cluster.ID)
	}

	err := kubeconfig.GenerateCombinedKubeconfig(c.rMSUrl, c.aPIToken, c.outputPath, clusterIDs)
	if err != nil {
		log.Fatalf("Run: error generating combined kubeconfig: %v", err)
	}
}
