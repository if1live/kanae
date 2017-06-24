package commands

import (
	"fmt"
	"os"
)

type Help struct {
}

func NewHelp() *Help {
	return &Help{}
}

const helpText = `
## supported command
* sync
* server
* dev
* help
`

func (cmd *Help) Execute() error {
	msgs := []string{
		fmt.Sprintf("Usage: %s --help\n", os.Args[0]),
	}
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	fmt.Println(helpText)
	return nil
}
