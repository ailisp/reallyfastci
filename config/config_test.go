package config

import "testing"

func TestLoadFullConfig(t *testing.T) {
	LoadConfig("../test/data/config/full.yaml")
}
