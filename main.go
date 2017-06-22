package main

import (
	"io/ioutil"
	"path"

	"github.com/if1live/kanae/settings"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	filename := "config.yaml"
	filepath := path.Join(GetExecutablePath(), filename)

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		check(err)
	}

	var s settings.Settings
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		check(err)
	}
}
