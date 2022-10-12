package layout

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/szpnygo/gtc/gocui"
	"github.com/szpnygo/gtc/log"
)

var data = `　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　◆　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　◆◆◆◆◆◆　　　　　　　　　　　　　　　　　　　◆◆◆　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　◆◆◆◆◆◆◆　　　　　　　　　　　　　　　　　　◆◆◆◆◆◆◆　　　　　　　　　　　　　　　　　◆◆◆◆◆◆　　　　　　　　　
　　　　　　　◆◆◆　　◆◆◆　　　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　　◆◆◆◆◆◆◆◆　　　　　　　　
　　　　　　　◆◆◆　　　◆◆　　　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　◆◆◆　　　◆◆◆　　　　　　　　
　　　　　　　◆◆◆　　　◆◆　　　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　◆◆◆　　　◆◆　　　　　　　　　
　　　　　　　◆◆◆　　　◆◆　　　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　
　　　　　　　◆◆◆　　◆◆◆　　　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　
　　　　　　　◆◆◆◆◆◆◆◆　　　　　　　　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　◆◆◆　　　　　◆　　　　　　　　
　　　　　　　◆◆◆◆◆◆◆◆　　　　　　　　　　　　　　　　　　　　◆◆　　　◆　　　　　　　　　　　　　　◆◆◆　　　　◆◆　　　　　　　　
　　　　　　　　　◆◆◆　◆◆　　　　　　　　　　　　　　　　　　　　◆◆◆◆◆◆　　　　　　　　　　　　　　　◆◆◆◆◆◆◆◆　　　　　　　　
　　　　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　　　◆◆◆◆◆　　　　　　　　　　　　　　　　◆◆◆◆◆◆◆　　　　　　　　　
　　　　　　　◆　　　　◆◆◆　　　　　　　　　　　　　　　　　　　　　◆◆◆　　　　　　　　　　　　　　　　　　　◆◆◆　　　　　　　　　　　
　　　　　　　◆◆◆◆◆◆◆◆　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　◆◆◆◆◆◆◆　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　
　　　　　　　　　　◆◆　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　`

type LayoutManager struct {
	g           *gocui.Gui
	width       int
	height      int
	currentView string
	currentRoom string
	username    string

	messages chan string

	onlineCount         int
	peerConnectionCount int

	stop            func(g *gocui.Gui, v *gocui.View) error
	onLoginEvent    func()
	onJoinRoomEvent func(string, string)
}

func NewLayoutManager(g *gocui.Gui, stop func(g *gocui.Gui, v *gocui.View) error) *LayoutManager {
	return &LayoutManager{
		g:        g,
		messages: make(chan string, 100),
		stop:     stop,
	}
}

func (lm *LayoutManager) OnLoginEvent(f func()) {
	lm.onLoginEvent = f
}

func (lm *LayoutManager) OnJoinRoomEvent(f func(string, string)) {
	lm.onJoinRoomEvent = f
}

func (lm *LayoutManager) GetMessage() chan string {
	return lm.messages
}

func (lm *LayoutManager) GetUsername() string {
	return lm.username
}

func (lm *LayoutManager) WriteMessage(name, msg string) {
	lm.g.Update(func(g *gocui.Gui) error {
		ov, err := g.View("messages")
		if err != nil {
			pterm.Error.Println("Cannot get output view:", err)
		}
		_, err = fmt.Fprintf(ov, "( %s  ): %s\n", name, msg)
		if err != nil {
			pterm.Error.Println("Cannot print to output view:", err)
		}
		ov.Rewind()
		log.GTCLog.Infoln("write message to view")

		return nil
	})
}

func (lm *LayoutManager) LoginLayout(g *gocui.Gui) error {
	width, height := g.Size()
	lm.width = width
	lm.height = height

	_ = lm.loginView(g)

	return nil
}

func (lm *LayoutManager) mainLayout(g *gocui.Gui) error {
	width, height := g.Size()
	lm.width = width
	lm.height = height

	_ = lm.messageBarLayout(g)
	_ = lm.infoView(g)
	_ = lm.roomsView(g)
	_ = lm.messagesView(g)
	_ = lm.chatView(g)
	_ = lm.usersView(g)

	return nil
}

func (lm *LayoutManager) messageBarLayout(g *gocui.Gui) error {
	v, err := g.SetView("bar", 0, 0, lm.width, 2)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	fmt.Fprintf(v, "GTC: ")
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	return nil
}

func (lm *LayoutManager) messagesView(g *gocui.Gui) error {
	v, err := g.SetView("messages", 24, 3, lm.width-32, lm.height-5)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	if len(lm.currentRoom) == 0 {
		v.Title = " Messages -- Please select a room from the left "
	} else {
		v.Title = fmt.Sprintf(" Messages -- %s ", lm.currentRoom)
	}

	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	v.Autoscroll = true
	return nil
}

func (lm *LayoutManager) infoView(g *gocui.Gui) error {
	v, err := g.SetView("info", 0, 3, 22, 9)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	v.Title = " Info "
	fmt.Fprintf(v, "Name: %s\n", lm.GetUsername())
	fmt.Fprintf(v, "Users: %d\n", lm.onlineCount)
	fmt.Fprintf(v, "PeerConnections: %d\n", lm.onlineCount)
	fmt.Fprintf(v, "\nTab to switch window\n")

	return nil
}

func (lm *LayoutManager) roomsView(g *gocui.Gui) error {
	v, err := g.SetView("rooms", 0, 10, 22, lm.height-2)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}

	v.Title = " 😁Rooms "
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	v.Autoscroll = false
	_, _ = g.SetCurrentView("rooms")
	lm.currentView = "rooms"

	return nil
}

func (lm *LayoutManager) usersView(g *gocui.Gui) error {
	v, err := g.SetView("users", lm.width-30, 3, lm.width-2, lm.height-2)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	v.Title = " Room Users "
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	_ = g.SetKeybinding("users", gocui.MouseWheelDown, gocui.ModNone, func(g *gocui.Gui, view *gocui.View) error {
		x, y := v.Origin()
		_ = v.SetOrigin(x, y+1)
		return nil
	})
	_ = g.SetKeybinding("users", gocui.MouseWheelUp, gocui.ModNone, func(g *gocui.Gui, view *gocui.View) error {
		x, y := v.Origin()
		_ = v.SetOrigin(x, y-1)
		return nil
	})

	return nil
}

func (lm *LayoutManager) chatView(g *gocui.Gui) error {
	v, err := g.SetView("chat", 24, lm.height-4, lm.width-32, lm.height-2)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	v.Title = " Press the enter to send the message "
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	v.Editable = true
	v.Wrap = true

	return nil
}

func (lm *LayoutManager) loginView(g *gocui.Gui) error {
	if _, err := g.SetView("login", 0, 0, lm.width-1, lm.height-1); err != nil {

	}

	if v, err := g.SetView("logo", lm.width/2-40, 2, lm.width/2+40, 30); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		_, err = fmt.Fprint(v, data)
		if err != nil {
			pterm.Error.Println("Cannot print to output view:", err)
		}
		v.Rewind()
	}

	if v, err := g.SetView("login_input", lm.width/2-30, 32, lm.width/2+30, 34); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Please input your name (Please press enter) "
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Editable = true
		v.Wrap = true

		_ = v.SetCursor(0, 0)
		_, _ = g.SetCurrentView("login_input")
	}

	return nil
}

func (lm *LayoutManager) KeyBinding(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, lm.stop); err != nil {
		pterm.Error.Println("Failed to set keybinding", err)
	}

	_ = g.SetKeybinding("login_input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		lm.username = v.Buffer()
		lm.username = strings.ReplaceAll(lm.username, "\n", "")
		if len(lm.username) == 0 {
			return nil
		}
		log.GTCLog.Infoln("key enter username", lm.username)
		g.SetManagerFunc(lm.mainLayout)
		lm.KeyBinding(g)
		lm.onLoginEvent()
		return nil
	})

	lm.ChatKeyBinding(g)
	lm.RoomKeyBinding(g)
}

func (lm *LayoutManager) ChatKeyBinding(g *gocui.Gui) {
	_ = g.SetKeybinding("chat", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(lm.currentRoom) == 0 {
			lm.UpdateMessageBar("Please select a chat room first", "red")
			return nil
		}
		v.Rewind()
		msg := v.Buffer()
		if len(msg) == 0 {
			return nil
		}
		if msg[len(msg)-1] == '\n' {
			msg = msg[:len(msg)-1]
		}
		lm.messages <- msg
		v.Clear()
		err := v.SetCursor(0, 0)
		if err != nil {
			pterm.Error.Println("Failed to set cursor:", err)
		}

		return nil
	})

	_ = g.SetKeybinding("chat", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if v == nil || v.Name() == "rooms" {
			return lm.SetInChatWindow(g)
		}
		return lm.SetInRoomWindow(g)
	})
	_ = g.SetKeybinding("chat", gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return lm.SetInChatWindow(g)
	})
}

func (lm *LayoutManager) RoomKeyBinding(g *gocui.Gui) {
	_ = g.SetKeybinding("rooms", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if v == nil || v.Name() == "rooms" {
			return lm.SetInChatWindow(g)
		}
		return lm.SetInRoomWindow(g)
	})
	_ = g.SetKeybinding("rooms", gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return nil
	})

	var up = func(g *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		_ = v.SetOrigin(x, y-1)
		return nil
	}
	var down = func(g *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		_ = v.SetOrigin(x, y+1)
		return nil
	}
	var getLine = func(g *gocui.Gui, v *gocui.View) error {
		var l string
		var err error

		_, cy := v.Cursor()
		if l, err = v.Line(cy); err != nil {
			l = ""
		}
		if len(l) > 0 {
			_ = lm.UpdateSelectedRoom(g, l)
		}

		return nil
	}

	_ = g.SetKeybinding("rooms", gocui.MouseWheelDown, gocui.ModNone, down)
	_ = g.SetKeybinding("rooms", gocui.MouseWheelUp, gocui.ModNone, up)
	_ = g.SetKeybinding("rooms", gocui.KeyArrowDown, gocui.ModNone, down)
	_ = g.SetKeybinding("rooms", gocui.KeyArrowUp, gocui.ModNone, up)
	_ = g.SetKeybinding("rooms", gocui.KeyEnter, gocui.ModNone, getLine)
	_ = g.SetKeybinding("rooms", gocui.MouseLeft, gocui.ModNone, getLine)

	lm.UpdateMessageBar("Welcome to gtc", "green")
}

func (lm *LayoutManager) SetInChatWindow(g *gocui.Gui) error {
	lm.currentView = "chat"
	v, err := g.SetCurrentView("chat")
	if err != nil {
		return err
	}
	v.Title = " 😁Press the enter to send the message "

	rv, err := g.View("rooms")
	if err != nil {
		return err
	}
	rv.Title = " Rooms "

	return nil
}

func (lm *LayoutManager) SetInRoomWindow(g *gocui.Gui) error {
	lm.currentView = "rooms"
	v, err := g.SetCurrentView("rooms")
	if err != nil {
		return err
	}
	v.Title = " 😁Rooms "

	rv, err := g.View("chat")
	if err != nil {
		return err
	}
	rv.Title = " Press the enter to send the message "

	return nil
}

func (lm *LayoutManager) UpdateSelectedRoom(g *gocui.Gui, room string) error {
	if lm.currentRoom == room {
		return nil
	}
	oldRoom := lm.currentRoom
	lm.currentRoom = room
	v, err := g.View("messages")
	if err != nil {
		return err
	}
	if len(lm.currentRoom) == 0 {
		v.Title = " Messages -- Please select a room from the left "
	} else {
		v.Title = fmt.Sprintf(" Messages -- %s ", lm.currentRoom)
	}
	v.Clear()

	if v, err := g.View("users"); err == nil {
		v.Clear()
	}

	lm.onJoinRoomEvent(oldRoom, lm.currentRoom)

	return nil
}

func (lm *LayoutManager) UpdateMessageBar(msg string, color string) {
	log.GTCLog.Warning(msg)
	colorNum := 31
	switch color {
	case "green":
		colorNum = 32
	case "red":
		colorNum = 31
	case "yellow":
		colorNum = 33
	case "blue":
		colorNum = 34
	}
	lm.g.Update(func(g *gocui.Gui) error {
		if v, err := g.View("bar"); err == nil {
			v.Clear()
			fmt.Fprintf(v, "GTC: \033[01;%dm%s\033[0m\n", colorNum, msg)
		}

		return nil
	})
}

func (lm *LayoutManager) UpdateOnlineCount(c int) {
	lm.onlineCount = c
}

func (lm *LayoutManager) UpdateRoomList(rooms []string) {
	lm.g.Update(func(g *gocui.Gui) error {
		v, err := g.View("rooms")
		if err != nil {
			return err
		}
		v.Clear()
		for _, name := range rooms {
			fmt.Fprintln(v, name)
		}
		return nil
	})
}

func (lm *LayoutManager) UpdateUserList(users []string) {
	lm.g.Update(func(g *gocui.Gui) error {
		v, err := g.View("users")
		if err != nil {
			return err
		}
		v.Clear()
		for _, name := range users {
			fmt.Fprintln(v, name)
		}
		return nil
	})
}
