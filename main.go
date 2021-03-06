package main

import (
	"path"
	"runtime"

	"github.com/if1live/kanae/commands"
	"github.com/if1live/kanae/kanaelib"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	filename := "config.yaml"
	filepath := path.Join(kanaelib.GetExecutablePath(), filename)
	s, err := kanaelib.LoadSettings(filepath)
	if err != nil {
		kanaelib.Check(err)
	}

	dispatcher := commands.NewDispatcher(s)
	dispatcher.Execute()
}
