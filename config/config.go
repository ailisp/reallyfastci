package config

import (
	"github.com/jinzhu/configor"
	"log"
)

var Config = struct {
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

func LoadConfig() {
	configor.Load(&Config, "config.yml")
	log.Printf("config: %#v", Config)
}
