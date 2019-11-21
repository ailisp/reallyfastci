package config

import (
	"github.com/jinzhu/configor"
	"log"
	"strings"
)

var Config = struct {
	ReallyfastciUrl     string   `yaml:"reallyfastci_url" required:"true"`
	GithubToken         string   `yaml:"github_token" required:"true" env:"GITHUB_TOKEN"`
	RepoUrl             string   `yaml:"repo_url" required:"true"`
	PushTriggerBranches []string `yaml:"push_trigger_branches"`

	Machine struct {
		MaxMachines  uint64 `yaml:"max_machines" default:"10"`
		MachineType  string `yaml:"machine_type" default:"n1-standard-32"`
		DiskSizeGB   uint64 `yaml:"disk_size_gb" default:"50"`
		ImageProject string `yaml:"image_project" default:"ubuntu-os-cloud"`
		ImageFamily  string `yaml:"image_family" default:"ubuntu-1804-lts"`
		Zone         string `default:" -west2-a"`
	}

	Build struct {
		Script         string `default:"./build.sh"`
		TimeoutMinutes uint64 `yaml:"timeout_minutes" default:"30"`
	}
}{}

var RepoName string

func LoadConfig(path string) {
	err := configor.Load(&Config, path)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	RepoName = getRepoName(Config.RepoUrl)
	log.Printf("config: %+v", Config)
}

func getRepoName(url string) (name string) {
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		log.Fatal("config error: invalid repo url")
	}
	if parts[len(parts)-1] == "" {
		name = strings.Join(parts[len(parts)-3:len(parts)-1], "/")
	} else {
		name = strings.Join(parts[len(parts)-2:len(parts)], "/")
	}
	return
}
