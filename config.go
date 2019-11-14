package pathfinder

import (
	"gopkg.in/yaml.v2"
	"log"
)

type Config struct {
	Pathfinder ConfigPathfinder `yaml:"pathfinder"`
}

type ConfigPathfinder struct {
	Path    ConfigPath              `yaml:"path"`
	TagMap  ConfigTagMap            `yaml:"tagmap"`
	Origins map[string]ConfigOrigin `yaml:"origins"`
}

type ConfigPath struct {
	Cue    string `yaml:"cue"`
	Prefix string `yaml:"prefix"`
	Suffix string `yaml:"suffix"`
}

type ConfigTagMap struct {
	Separator string      `yaml:"separator"`
	Trim      string      `yaml:"trim"`
	Tags      []ConfigTag `yaml:"tags"`
}

type ConfigTag struct {
	Tag   string `yaml:"tag"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type ConfigOrigin struct {
	Enable bool `yaml:"bool"`
}

func NewConfig() *Config {
	var c Config

	return &c
}

func ReadConfig(data []byte) (*Config, error) {
	cnf := NewConfig()

	err := yaml.Unmarshal(data, cnf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return cnf, err
}
