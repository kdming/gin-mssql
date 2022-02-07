package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type config struct {
	SqlUrl  string `yaml:"SqlUrl"`
	JwtKey  string `yaml:"JwtKey"`
	AppPort string `yaml:"AppPort"`
}

func GetConfig() *config {
	root, err := os.Getwd()
	filePath := root + "/config.yaml"
	buffer, err := ioutil.ReadFile(filePath)
	conf := &config{}
	err = yaml.Unmarshal(buffer, conf)
	if err != nil {
		panic(err)
	}
	return conf
}
