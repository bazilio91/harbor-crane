package crane

import (
	"net/http"
	"strings"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/sirupsen/logrus"
)

func (c *Crane) GetRegistry(url string) (*registry.Registry, error) {
	if url == "docker.io" {
		url = "registry-1.docker.io"
	}

	crStore := c.Config.DockerConfig.GetCredentialsStore(url)
	creds, err := crStore.Get(url)
	if err != nil {
		return nil, err
	}

	repo_url := url
	if !strings.HasPrefix(url, "http") {
		repo_url = "https://" + url
	}

	reg := createRegistry(repo_url, creds.Username, creds.Password, logrus.Debugf)

	return reg, nil
}

func createRegistry(registryURL, username, password string, logf registry.LogfCallback) *registry.Registry {
	transport := http.DefaultTransport

	url := strings.TrimSuffix(registryURL, "/")
	transport = registry.WrapTransport(transport, url, username, password)
	registry := &registry.Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
		Logf: logf,
	}
	return registry
}
