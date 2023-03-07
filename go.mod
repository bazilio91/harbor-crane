module github.com/bazilio91/harbor-crane

require (
	github.com/avast/retry-go v2.4.1+incompatible
	github.com/docker/cli v0.0.0-20190822175708-578ab52ece34
	github.com/docker/distribution v2.7.1+incompatible
	github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/yaml.v2 v2.2.2
)

require (
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.0.0-rc8 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5 => github.com/kpaas-io/docker-registry-client v0.0.0-20190926031228-4141b50ff9fa
