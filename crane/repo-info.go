package crane

import (
	"github.com/docker/distribution/reference"
	"github.com/heroku/docker-registry-client/registry"
)

func (c *Crane) GetRepoInfo(url string, reg *registry.Registry) (reference.Named, error) {
	ref, err := reference.ParseNormalizedNamed(url)
	if err != nil {
		return nil, err
	}

	return ref, nil
}
