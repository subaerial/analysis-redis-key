package config

import (
	"analysis.redis/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var Properties *model.Properties

func InitProperties() {
	file, errf := ioutil.ReadFile("application.yml")
	if errf != nil {
		log.Fatal("fail to read file:", errf)
	}
	err := yaml.Unmarshal(file, &Properties)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
