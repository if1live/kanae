package commands

import "github.com/if1live/kanae/kanaelib"

type Server struct {
	settings kanaelib.Settings
}

func NewServer(s kanaelib.Settings) *Server {
	return &Server{
		settings: s,
	}
}

func (cmd *Server) Execute() error {
	return nil
}
