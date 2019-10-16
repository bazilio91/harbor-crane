package crane

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/docker/distribution/reference"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

func (c *Crane) TransferBlob(
	sourceRegistry *registry.Registry, destRegistry *registry.Registry,
	sourceReference reference.NamedTagged, destReference reference.NamedTagged,
	d digest.Digest) error {
	return retry.Do(
		func() error {
			exists, err := destRegistry.HasBlob(reference.Path(sourceReference), d)
			if exists {
				return nil
			}

			logrus.Infof("uploading blob %s", d)
			reader, err := sourceRegistry.DownloadBlob(reference.Path(sourceReference), d)
			if reader != nil {
				defer reader.Close()
			}
			if err != nil {
				return err
			}

			err = destRegistry.UploadBlob(reference.Path(destReference), d, reader)

			if err != nil {
				logrus.WithError(err).Error("failed to upload blob")
				return err
			}

			return nil
		},
		retry.OnRetry(func(n uint, err error) {
			logrus.Warnf("Blob push failed, retry %v", n+1)
		}),
		retry.Delay(1*time.Second),
		retry.Attempts(1),
	)
}
