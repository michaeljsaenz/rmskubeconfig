package types

import "fmt"

type RMSCluster struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RMSClusterResponse struct {
	Data []RMSCluster `json:"data"`
}

type KubeconfigResponse struct {
	Config string `json:"config"`
}

type KubeconfigClusterDetails struct {
	Server                   string `yaml:"server" json:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data" json:"certificate-authority-data"`
}

type KubeconfigCluster struct {
	Name    string                   `yaml:"name" json:"name"`
	Cluster KubeconfigClusterDetails `yaml:"cluster" json:"cluster"`
}

type KubeconfigUser struct {
	Name string `yaml:"name" json:"name"`
	User struct {
		Token string `yaml:"token" json:"token"`
	} `yaml:"user" json:"user"`
}

type KubeconfigContext struct {
	Name    string `yaml:"name" json:"name"`
	Context struct {
		User    string `yaml:"user" json:"user"`
		Cluster string `yaml:"cluster" json:"cluster"`
	} `yaml:"context" json:"context"`
}

type Kubeconfig struct {
	APIVersion string              `yaml:"apiVersion" json:"apiVersion"`
	Kind       string              `yaml:"kind" json:"kind"`
	Clusters   []KubeconfigCluster `yaml:"clusters" json:"clusters"`
	Users      []KubeconfigUser    `yaml:"users" json:"users"`
	Contexts   []KubeconfigContext `yaml:"contexts" json:"contexts"`
}

const ErrRequestCode = 1000

type RequestError struct {
	Code    int
	Message string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
