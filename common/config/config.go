package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	SqlUrl  string `yaml:"SqlUrl"`
	JwtKey  string `yaml:"JwtKey"`
	AppPort string `yaml:"AppPort"`
}

func GetConfig() *Config {
	root, err := os.Getwd()
	filePath := root + "/config.yaml"
	buffer, err := ioutil.ReadFile(filePath)
	conf := &Config{}
	err = yaml.Unmarshal(buffer, conf)
	if err != nil {
		panic(err)
	}
	return conf
}
