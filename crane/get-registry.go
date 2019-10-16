package crane

import (
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

	reg, err := registry.New(repo_url, creds.Username, creds.Password)

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
