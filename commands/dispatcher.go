package commands

import (
	"fmt"

	"flag"

	"os"

	"github.com/if1live/kanae/kanaelib"
)

var command string

// for server
var port int
var host string

func init() {
	flag.StringVar(&command, "command", "", "command, use help command to see detail")
	flag.IntVar(&port, "port", 8000, "port to use")
	flag.StringVar(&host, "host", "127.0.0.1", "address to use")
}

type Dispatcher struct {
	settings kanaelib.Settings
	commands map[string]Command
}

func NewDispatcher(s kanaelib.Settings) Dispatcher {
	cmds := map[string]Command{
		"server": NewServer(s),
		"sync":   NewSync(s),
		"dev":    NewDev(s),
		"help":   NewHelp(),
	}
	return Dispatcher{
		settings: s,
		commands: cmds,
	}
}

func (d *Dispatcher) getCommand(command string) Command {
	if command == "" {
		return d.commands["help"]
	}

	cmd, ok := d.commands[command]
	if !ok {
		return nil
	}
	return cmd
}

func (d *Dispatcher) Execute() {
	flag.Parse()

	cmd := d.getCommand(command)
	if cmd == nil {
		fmt.Println("invalid command :", command)
		os.Exit(-1)
	}

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
