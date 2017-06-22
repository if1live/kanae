package main

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

func loadConfig(filepath string) (Settings, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Settings{}, err
	}

	var s Settings
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		return Settings{}, err
	}

	return s, nil
}

func main() {
	filename := "config.yaml"
	filepath := path.Join(GetExecutablePath(), filename)
	s, err := loadConfig(filepath)
	if err != nil {
		check(err)
	}

	report := NewPoloniexReport(s)
	report.Generate()
}
