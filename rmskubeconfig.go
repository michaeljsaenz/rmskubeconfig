package rmskubeconfig

import (
	"errors"
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
func NewConfig() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Config{
		rMSUrl:     "",
		aPIToken:   "",
		outputPath: cwd,
	}, nil
}

// SetRMSUrl sets RMS API URL
func (c *Config) SetRMSUrl(url string) error {
	// validate URL format
	regex := `^(https?:\/\/)?([\w\-]+(\.[\w\-]+)+)(:[0-9]{1,5})?(\/[^\s]*)?$`
	if match, _ := regexp.MatchString(regex, url); !match {
		return errors.New("invalid RMS URL format")
	}
	c.rMSUrl = url
	return nil
}

// SetApiToken sets RMS API token
func (c *Config) SetApiToken(token string) error {
	// validate token format
	regex := `^token-\w+:\w+`

	if match, _ := regexp.MatchString(regex, token); !match {
		return errors.New("invalid API token format")
	}
	c.aPIToken = token
	return nil
}

// SetOutputPath sets path where to save config file
func (c *Config) SetOutputPath(path string) error {
	// validate path exists and ensure it's a directory
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		return errors.New("output path must be an existing directory")
	}
	// convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return errors.New("failed to resolve absolute path")
	}
	c.outputPath = absPath
	return nil
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

	err := kubeconfig.GenerateCombinedKubeconfig(c.rMSUrl, c.aPIToken, clusterIDs)
	if err != nil {
		log.Fatalf("error generating combined kubeconfig: %v", err)
	}
}
