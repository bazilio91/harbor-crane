package crane

import (
	"fmt"
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

	reg, err := createRegistry(repo_url, creds.Username, creds.Password, logrus.Debugf)

	if err != nil {
		logrus.WithError(err).Error("failed to create https repo, trying http")

		repo_url = url

		if !strings.HasPrefix(url, "http") {
			repo_url = "http://" + url
		}

		reg, err = registry.New(repo_url, creds.Username, creds.Password)

		if err != nil {
			return nil, err
		}
	}

	reg.Logf = logrus.Debugf

	return reg, nil
}

func createRegistry(registryURL, username, password string, logf registry.LogfCallback) (*registry.Registry, error) {
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

	resp, err := registry.Client.Get(fmt.Sprintf("%s%s", registry.URL, "/v2/_ping"))
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp != nil && resp.StatusCode < 500 {
		return registry, nil
	}

	return nil, err
}
