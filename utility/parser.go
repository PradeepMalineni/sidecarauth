package logger

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Server struct {
	IP string `yaml:"ip"`
}

type Application struct {
	Name    string   `yaml:"name"`
	Servers []Server `yaml:"servers"`
}

type Config struct {
	Application Application `yaml:"application"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
