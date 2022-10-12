package client

import (
	"errors"

	"github.com/pterm/pterm"
	"github.com/szpnygo/gtc/client/layout"
	"github.com/szpnygo/gtc/client/server"
	"github.com/szpnygo/gtc/gocui"
)

type GTCClient struct {
	g   *gocui.Gui
	l   *layout.LayoutManager
	s   *server.ClientServer
	api string
}

func NewGTCClient(api string) *GTCClient {
	return &GTCClient{
		api: api,
	}
}

func (gtc *GTCClient) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		pterm.Error.Println("Failed to create GUI", err)
		return err
	}
	defer g.Close()
	gtc.g = g

	//config
	g.Mouse = false
	g.Cursor = true
	g.InputEsc = true

	gtc.l = layout.NewLayoutManager(g, func(g *gocui.Gui, v *gocui.View) error {
		if gtc.s != nil {
			gtc.s.Stop()
		}
		return gocui.ErrQuit
	})

	gtc.s = server.NewClientServer(gtc.api, gtc.l)
	go gtc.s.Run()

	g.SetManagerFunc(gtc.l.LoginLayout)
	gtc.l.KeyBinding(g)

	if err := g.MainLoop(); err != nil && errors.Is(err, gocui.ErrQuit) {
		pterm.Error.Println("Failed to start main loop", err)
	}

	return nil
}
