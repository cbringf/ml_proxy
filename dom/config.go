package dom

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MLApi struct {
		Host string `yaml:"host"`
	} `yaml:"ml_api"`
	DB struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"db"`
}

func Load() (*Config, error) {
	f, err := os.Open("config.yaml")

	if err != nil {
		log.Fatal("Unable to load initial config")

		return nil, err
	}

	defer f.Close()

	var result Config

	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&result)

	if err != nil {
		log.Fatal("Unable to decode initial config")

		return nil, err
	}

	return &result, nil
}
