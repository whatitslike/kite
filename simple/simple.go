package simple

import (
	"github.com/koding/kite"
	"github.com/koding/kite/config"
	"github.com/koding/kite/kontrolclient"
	"github.com/koding/kite/registration"
	"github.com/koding/kite/server"
)

// simple kite server
type Simple struct {
	*server.Server
	Kontrol      *kontrolclient.Kontrol
	Registration *registration.Registration
}

func New(name, version string) *Simple {
	config := config.New()

	err := config.ReadKiteKey()
	if err != nil {
		panic(err)
	}
	config.ReadEnvironmentVariables()

	k := kite.New(name, version)
	k.Config = config

	server := server.New(k)

	kon := kontrolclient.New(k)

	s := &Simple{
		Server:       server,
		Kontrol:      kon,
		Registration: registration.New(kon),
	}

	return s
}

// HandleFunc registers a handler to run when a method call is received from a Kite.
func (s *Simple) HandleFunc(method string, handler kite.HandlerFunc) {
	s.Server.Kite.HandleFunc(method, handler)
}

func (s *Simple) Run() {
	s.Kontrol.DialForever()
	s.Server.Start()
	s.Registration.RegisterToProxyAndKontrol()
	<-s.Server.CloseNotify()
}

func (s *Simple) Close() {
	s.Server.Close()
}