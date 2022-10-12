package server

import (
	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/config"
	"github.com/szpnygo/gtc/server/internal/server"
	"github.com/szpnygo/gtc/server/internal/svc"
	"github.com/szpnygo/gtc/server/model"
)

func Server(rooms string) {

	config := config.NewDefaultConfig()
	// config.Debug = true

	app := gogs.NewApp(config)
	app.AddAcceptor(acceptor.NewWSAcceptror(&acceptor.AcceptroConfig{
		HttpPort: 8888,
		Name:     "gtc",
		Groups: []*acceptor.AcceptorGroupConfig{
			&acceptor.AcceptorGroupConfig{
				GroupName: "gtc",
			},
		},
	}))

	ctx := svc.NewServiceContext(app, rooms)
	srv := server.NewServer(ctx)

	model.RegisterAllComponents(app, srv)

	defer app.Shutdown()
	app.Start()
}
