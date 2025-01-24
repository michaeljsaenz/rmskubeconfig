package rmskubeconfig

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/michaeljsaenz/rmskubeconfig/internal/kubeconfig"
)

// Config holds values for processing
type Config struct {
	rmsUrl     string
	apiToken   string
	outputPath string
}

// NewConfig creates a new Config instance with default values
func NewConfig() *Config {
	return &Config{
		rmsUrl:     "",
		apiToken:   "",
		outputPath: "",
	}
}

// SetRMSUrl sets RMS API URL
func (c *Config) SetRMSUrl(url string) error {
	// validate URL format
	regex := `^(https?:\/\/)?([\w\-]+(\.[\w\-]+)+)(:[0-9]{1,5})?(\/[^\s]*)?$`
	if match, _ := regexp.MatchString(regex, url); !match {
		return fmt.Errorf("invalid RMS URL format: %s", url)
	}
	c.rmsUrl = url
	return nil
}

// SetApiToken sets RMS API token
func (c *Config) SetApiToken(token string) error {
	// validate token format
	regex := `^token-\w+:\w+`

	if match, _ := regexp.MatchString(regex, token); !match {
		return fmt.Errorf("invalid API token format, must match regex: %q", regex)
	}
	c.apiToken = token
	return nil
}

// SetOutputPath sets path where to save config file
func (c *Config) SetOutputPath(path string) error {
	// validate path exists and ensure it's a directory
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		return fmt.Errorf("output path must be an existing directory: %s", path)
	}

	c.outputPath = path

	return nil
}

// RMSUrl returns RMS API URL
func (c *Config) RMSUrl() string {
	return c.rmsUrl
}

// ApiToken returns RMS API token
func (c *Config) ApiToken() string {
	return c.apiToken
}

// OutputPath returns output file path
func (c *Config) OutputPath() string {
	return c.outputPath
}

// Run executes the Config to generate combined kubeconfig (config) file
func (c *Config) Run() {

	if c.outputPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current working directory: %v", err)
		}

		c.outputPath = cwd
	}

	// convert to absolute path
	absPath, err := filepath.Abs(c.outputPath)
	if err != nil {
		log.Fatalf("failed to resolve absolute path: %s, error: %v", c.outputPath, err)
	}

	c.outputPath = absPath

	clusters, _ := kubeconfig.GetClusters(c.rmsUrl, c.apiToken)

	var clusterIDs []string
	for _, cluster := range clusters {
		clusterIDs = append(clusterIDs, cluster.ID)
	}

	err = kubeconfig.GenerateCombinedKubeconfig(c.rmsUrl, c.apiToken, c.outputPath, clusterIDs)
	if err != nil {
		log.Fatalf("Run: error generating combined kubeconfig: %v", err)
	}
}
