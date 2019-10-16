module github.com/bazilio91/harbor-crane

require (
	github.com/avast/retry-go v2.4.1+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/cli v0.0.0-20190822175708-578ab52ece34
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc8 // indirect
	github.com/ryanuber/go-glob v0.0.0-20170128012129-256dc444b735
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5 => github.com/kpaas-io/docker-registry-client v0.0.0-20190926031228-4141b50ff9fa
