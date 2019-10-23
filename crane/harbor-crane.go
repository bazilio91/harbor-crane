package crane

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/bazilio91/harbor-crane/config"
	"github.com/docker/distribution/reference"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

type Crane struct {
	Source *registry.Registry
	Dest   *registry.Registry
	Config config.CraneConfig
}

func NewCrane() *Crane {
	sync := new(Crane)

	sync.Config = config.NewConfig()

	logrus.SetLevel(logrus.InfoLevel)

	return sync
}

func (c *Crane) SyncRepo(repoConfig config.RepoConfig, errors *[]error) {
	sourceReference, err := c.GetRepoInfo(repoConfig.Source, nil)
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	sourceRegistry, err := c.GetRegistry(reference.Domain(sourceReference))
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	destReference, err := c.GetRepoInfo(repoConfig.Dest, nil)
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	destRegistry, err := c.GetRegistry(reference.Domain(destReference))
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	sourceTags, err := sourceRegistry.Tags(reference.Path(sourceReference))
	if err != nil {
		*errors = append(*errors, err)
		return
	}

	tagsForSync := []string{}
	for _, tag := range repoConfig.Tags {
		for _, sourceTag := range sourceTags {
			if match, _ := regexp.MatchString(tag, sourceTag); match {
				tagsForSync = append(tagsForSync, sourceTag)
			}
		}
	}

	fmt.Printf("Tags for sync: %s\n", tagsForSync)

	for _, tag := range tagsForSync {
		destTagRef, err := reference.WithTag(destReference, tag)
		if err != nil {
			*errors = append(*errors, err)
			continue
		}

		sourceTagRef, err := reference.WithTag(sourceReference, tag)
		if err != nil {
			*errors = append(*errors, err)
			continue
		}

		err = c.SyncTag(sourceRegistry, destRegistry, sourceTagRef, destTagRef)
		if err != nil {
			*errors = append(*errors, err)
			continue
		}
	}
}

func (c *Crane) GetRepos(reg *registry.Registry) ([]string, error) {
	repos, err := reg.Repositories()

	if err != nil {
		return nil, err
	}
	b := repos[:0]
	for _, x := range repos {
		b = append(b, x)
	}

	return b, nil
}

func (c *Crane) GetTags(repo reference.Named, reg *registry.Registry) (repoTags map[string]digest.Digest, err error) {
	tags, err := reg.Tags(reference.Path(repo))
	if err != nil {
		return nil, err
	}

	repoTags = map[string]digest.Digest{}

	var wg sync.WaitGroup
	for _, t := range tags {
		wg.Add(1)
		go func(tag string) {
			defer wg.Done()

			d, err := reg.ManifestDigest(reference.Path(repo), tag)
			if err != nil {
				logrus.WithError(err).Error("failed to get digest")
				return
			}

			repoTags[t] = d
		}(t)
	}

	wg.Wait()

	return
}

func (c *Crane) SyncTag(sourceRegistry *registry.Registry, destRegistry *registry.Registry, sourceReference reference.NamedTagged,
	destReference reference.NamedTagged) error {
	logrus.Infof("syncing %s into %s", sourceReference.String(), destReference.String())

	destManifest, _ := destRegistry.ManifestV2(reference.Path(destReference), destReference.Tag())

	sourceManifest, err := sourceRegistry.ManifestV2(reference.Path(sourceReference), sourceReference.Tag())
	if err != nil {
		return err
	}

	if destManifest != nil && destManifest.Config.Digest.String() == sourceManifest.Config.Digest.String() {
		logrus.Infof("%s already exists", destReference.String())
		return nil
	}

	var wg = sync.WaitGroup{}

	logrus.Infof("%s is syncing", destReference.String())
	for _, l := range sourceManifest.Layers {
		wg.Add(1)
		d := l.Digest
		go func(layer digest.Digest) {
			defer wg.Done()

			err = c.TransferBlob(sourceRegistry, destRegistry, sourceReference, destReference, d)

			if err != nil {
				panic(err)
			}
		}(d)
	}

	wg.Wait()

	err = c.TransferBlob(sourceRegistry, destRegistry, sourceReference, destReference, sourceManifest.Config.Digest)
	if err != nil {
		return err
	}

	logrus.Infof("pushing manifest %s", sourceManifest.Config.Digest)
	err = destRegistry.PutManifest(reference.Path(destReference), destReference.Tag(), sourceManifest)

	return err
}

func (c *Crane) Sync() error {
	hadErr := false
	for _, repoConfig := range c.Config.Repos {
		repoErrors := make([]error, 0)

		c.SyncRepo(repoConfig, &repoErrors)
		if len(repoErrors) > 0 {
			hadErr = true
			logrus.Errorf("Repo %s synced with repoErrors:\n", repoConfig.Source)
			for _, e := range repoErrors {
				logrus.Errorln(e)
			}
		}
	}

	if hadErr {
		return errors.New("some repos failed to sync")
	}
	return nil
}
