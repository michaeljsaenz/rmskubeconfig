<p align="center">
  <a href="https://goreportcard.com/report/github.com/michaeljsaenz/rmskubeconfig"><img src="https://goreportcard.com/badge/github.com/michaeljsaenz/rmskubeconfig" alt="Code Status" ></a>
  <a href="https://codecov.io/github/michaeljsaenz/rmskubeconfig"><img src="https://codecov.io/github/michaeljsaenz/rmskubeconfig/graph/badge.svg?token=IYNU53BPM7"/></a>
  <a href="https://img.shields.io/github/v/release/michaeljsaenz/rmskubeconfig?include_prereleases" title="Latest Release" rel="nofollow"><img src="https://img.shields.io/github/v/release/michaeljsaenz/rmskubeconfig?include_prereleases" alt="Latest Release"></a>
</p>

# `rmskubeconfig`

`rmskubeconfig` aggregates multiple kubeconfig files from **RMS (Rancher Management Service) managed clusters** into a single configuration file for simplified access.

#### Table of Contents  
- [Features](#features)
- [Usage](#usage)
  - [Initialize Configuration](#initialize-configuration)
  - [Set RMS API URL](#set-rms-api-url)
  - [Set API Token](#set-api-token)
  - [Set Output Path](#set-output-path)
  - [Set Cluster ID (for scoped tokens)](#set-cluster-id-for-scoped-tokens)
  - [Generate Combined Kubeconfig](#generate-combined-kubeconfig)
- [Sample Package Use](#sample-package-use)
- [Usage with Scoped Tokens](#usage-with-scoped-tokens)


## Features
- **Configuration Management:** Stores RMS API URL, API token, and output path.
- **Input Validation:** Ensures RMS URL, API token, and output path are valid.
- **Cluster Retrieval:** Fetches kubeconfig of all RMS-managed clusters via the RMS API.
- **Scoped Token Support:** Works with RMS tokens that are scoped to specific cluster IDs.
- **Kubeconfig Generation:** Merges kubeconfig files into a unified configuration.

## Usage

### Initialize Configuration
```go
config := rmskubeconfig.NewConfig()
```

### Set RMS API URL
```go
err := config.SetRMSUrl("https://your-rms-api-url.com")
if err != nil {
    // handle error
}
```

### Set API Token
```go
err := config.SetApiToken("your-api-token")
if err != nil {
    // handle error
}
```

### Set Output Path
```go
err := config.SetOutputPath("/path/to/save/kubeconfig") // defaults to current-working-directory
if err != nil {
    // handle error
}
```

### Set Cluster ID (for scoped tokens)
```go
// Use this when your RMS token is scoped to a specific cluster and cannot list all clusters
err := config.SetClusterID("your-cluster-id")
if err != nil {
    // handle error
}
```

### Generate Combined Kubeconfig
```go
err := config.Run()
if err != nil {
    // handle error
}
```

## Sample Package Use
```go
package main

import (
	"log"
	"os"

	"github.com/michaeljsaenz/rmskubeconfig"
)

func main() {
	cfg := rmskubeconfig.NewConfig()
	cfg.SetApiToken(getEnv("RMS_TOKEN"))
	cfg.SetRMSUrl(getEnv("RMS_URL"))
	cfg.Run()
}

func getEnv(envKey string) (value string) {
	value, ok := os.LookupEnv(envKey)
	if !ok {
		log.Fatalf("Error: `%v` environment variable not set, must set.", envKey)
	}
	return
}

```

This generates a single kubeconfig file at the specified output path (current working directory by default).

## Usage with Scoped Tokens

If your RMS token is scoped to a specific cluster ID (i.e., it cannot list all clusters but can generate kubeconfig for a specific cluster), use the following approach:

```go
package main

import (
	"log"
	"os"

	"github.com/michaeljsaenz/rmskubeconfig"
)

func main() {
	cfg := rmskubeconfig.NewConfig()
	cfg.SetApiToken(getEnv("RMS_TOKEN"))
	cfg.SetRMSUrl(getEnv("RMS_URL"))
	cfg.SetClusterID(getEnv("CLUSTER_ID")) // Set the specific cluster ID for scoped tokens
	cfg.Run()
}

func getEnv(envKey string) (value string) {
	value, ok := os.LookupEnv(envKey)
	if !ok {
		log.Fatalf("Error: `%v` environment variable not set, must set.", envKey)
	}
	return
}

```

This approach bypasses the need to list all clusters and directly generates the kubeconfig for the specified cluster ID.
