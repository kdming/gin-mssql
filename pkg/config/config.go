package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	DB_HOST   string `yaml:"DB_HOST"`
	DB_USER   string `yaml:"DB_USER"`
	DB_PWD    string `yaml:"DB_PWD"`
	DB_NAME   string `yaml:"DB_NAME"`
	Token_KEY string `yaml:"Token_KEY"`
	HttpUrl   string `yaml:"HttpUrl"`
	TB_SYNC   bool   `yaml:"TB_SYNC"`
	APP_PORT  string `yaml:"APP_PORT"`
}

func GetConfig() *config {
	// 设置文件路径
	root, err := os.Getwd()
	filePath := root + "/config.yaml"

	// 读取并解析文件
	buffer, err := ioutil.ReadFile(filePath)
	config := &config{}
	err = yaml.Unmarshal(buffer, &config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return config
}
