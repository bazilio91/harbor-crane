package config

import (
	"io/ioutil"

	cliconfig "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"gopkg.in/yaml.v2"
)

type RepoConfig struct {
	Source string
	Dest   string
	Tags   []string
}

type CraneConfig struct {
	DockerConfig *configfile.ConfigFile

	Repos []RepoConfig
	DefaultRegistry string
}

func NewConfig() CraneConfig {
	c := CraneConfig{
		DockerConfig: cliconfig.LoadDefaultConfigFile(nil),
		DefaultRegistry: "https://index.docker.io",
	}

	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}

	return c
}
