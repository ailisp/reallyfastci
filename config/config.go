package config

import (
	"github.com/jinzhu/configor"
	"log"
	"strings"
)

var Config = struct {
	ReallyFastCiUrl     string `required:"true"`
	GithubToken         string `required:"true" env:"GITHUB_TOKEN"`
	RepoUrl             string `required:"true"`
	PushTriggerBranches []string

	Machine struct {
		MaxMachines  uint64 `default:"10"`
		MachineType  string `default:"n1-standard-32"`
		DiskSizeGB   uint64 `default:"50"`
		ImageProject string `default:"ubuntu-os-cloud"`
		ImageFamily  string `default:"ubuntu-1804-lts"`
		Zone         string `default:"us-west2-a"`
	}

	Build struct {
		Script         string `default:"./build.sh"`
		TimeoutMinutes uint64 `default:"30"`
	}
}{}

var RepoName string

func LoadConfig() {
	configor.Load(&Config, "config.yml")
	RepoName = getRepoName(Config.RepoUrl)
	log.Printf("config: %#v", Config)
}

func getRepoName(url string) (name string) {
	parts := strings.Split(url, "/")
	if parts[len(parts)-1] == "" {
		name = strings.Join(parts[len(parts)-3:len(parts)-1], "/")
	} else {
		name = strings.Join(parts[len(parts)-2:len(parts)], "/")
	}
	return
}
