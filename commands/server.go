package commands

import (
	"github.com/if1live/kanae/kanaelib"
	"github.com/if1live/kanae/web"
)

type Server struct {
	settings kanaelib.Settings
}

func NewServer(s kanaelib.Settings) *Server {
	return &Server{
		settings: s,
	}
}

func (cmd *Server) Execute() error {
	s := web.NewServer(host, port, cmd.settings)
	s.Run()
	return nil
}
