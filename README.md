<p align="center">
  <a href="https://goreportcard.com/report/github.com/michaeljsaenz/rmskubeconfig"><img src="https://goreportcard.com/badge/github.com/michaeljsaenz/rmskubeconfig" alt="Code Status" ></a>
  <a href="https://codecov.io/github/michaeljsaenz/rmskubeconfig"><img src="https://codecov.io/github/michaeljsaenz/rmskubeconfig/graph/badge.svg?token=IYNU53BPM7"/></a>
  <a href="https://img.shields.io/github/v/release/michaeljsaenz/rmskubeconfig?include_prereleases" title="Latest Release" rel="nofollow"><img src="https://img.shields.io/github/v/release/michaeljsaenz/rmskubeconfig?include_prereleases" alt="Latest Release"></a>
</p>

# `rmskubeconfig`

`rmskubeconfig` aggregates multiple kubeconfig files from **RMS (Rancher Management Service) managed clusters** into a single configuration file for simplified access.

## Features
- **Configuration Management:** Stores RMS API URL, API token, and output path.
- **Input Validation:** Ensures RMS URL, API token, and output path are valid.
- **Cluster Retrieval:** Fetches kubeconfig of all RMS-managed clusters via the RMS API.
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

### Generate Combined Kubeconfig
```go
err := config.Run()
if err != nil {
    // handle error
}
```

This generates a single kubeconfig file at the specified output path.
