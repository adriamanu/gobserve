package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Command      string   `json:"command" yaml:"command"`
	Files        []string `json:"files" yaml:"files"`
	IgnoredFiles []string `json:"ignoredFiles" yaml:"ignoredFiles"`
}

func ParseConfigFile(fp string) Config {
	ext := filepath.Ext(fp)
	if ext != ".json" && ext != ".yaml" && ext != ".yml" {
		log.Fatal("Configuration file must be in json or yaml format.")
	}

	file, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var conf Config
	if ext == ".json" {
		if err = json.Unmarshal(fileContent, &conf); err != nil {
			log.Fatal("error while parsing json", err)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		if err = yaml.Unmarshal(fileContent, &conf); err != nil {
			log.Fatal("error while parsing yaml", err)
		}
	}

	if conf.Command == "" {
		log.Fatal("command is a string and must be provided in the configuration file")
	}
	if len(conf.Files) <= 0 {
		log.Fatal("files is an array of string and must be provided in the configuration file")
	}
	return conf
}
