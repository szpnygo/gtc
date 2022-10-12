package server

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/codec"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/dispatch"
	"github.com/metagogs/gogs/global"
	"github.com/metagogs/gogs/proto"
	"github.com/pterm/pterm"
	"github.com/szpnygo/gtc/client/layout"
	"github.com/szpnygo/gtc/p2p"
	"github.com/szpnygo/gtc/server/model"

	gproto "google.golang.org/protobuf/proto"
)

type ClientServer struct {
	layout        *layout.LayoutManager
	done          chan struct{}
	dp            *dispatch.DispatchServer
	codecHelper   *codec.CodecHelper
	signalingConn *websocket.Conn
	clients       map[int64]*p2p.P2PClient
	api           string
}

func NewClientServer(api string, layout *layout.LayoutManager) *ClientServer {
	global.GOGS_DISABLE_LOG = true
	dp := dispatch.NewDispatchServer()
	codecHelper := codec.NewCodecHelper(&config.Config{}, dp)

	return &ClientServer{
		dp:          dp,
		codecHelper: codecHelper,
		layout:      layout,
		done:        make(chan struct{}),
		clients:     make(map[int64]*p2p.P2PClient),
		api:         api,
	}
}

func (c *ClientServer) Run() {
	c.dp.RegisterComponent(_GTCComponentDesc, c)
	c.layout.OnLoginEvent(c.login)
	c.layout.OnJoinRoomEvent(c.joinRoom)

	writeMessage := c.layout.GetMessage()
	for {
		select {
		case msg := <-writeMessage:
			c.SendMessage(msg)
		case <-c.done:
			return
		}
	}
}

func (c *ClientServer) SendMessage(msg string) {
	c.layout.WriteMessage(c.layout.GetUsername(), msg)
	message := &model.Message{
		Name: c.layout.GetUsername(),
		Data: msg,
	}
	if data, err := gproto.Marshal(message); err == nil {
		for _, c := range c.clients {
			c.SendMessage(data)
		}
	}
}

func (c *ClientServer) Stop() {
	close(c.done)
}

func (c *ClientServer) login() {
	// connect to signaling server
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/gtc", c.api), nil)
	if err != nil {
		pterm.Error.Println(err)
		panic(err)
	}
	c.signalingConn = conn
	go c.readSignalingMessage()

	c.sendSignalingMessage(&model.ListRoom{})
}

func (c *ClientServer) joinRoom(leave, join string) {
	leaveRoom := &model.LeaveRoom{
		RoomId: leave,
		Name:   c.layout.GetUsername(),
	}
	c.sendSignalingMessage(leaveRoom)
	joinRoom := &model.JoinRoom{
		RoomId: join,
		Name:   c.layout.GetUsername(),
	}
	c.sendSignalingMessage(joinRoom)
}

func (c *ClientServer) readSignalingMessage() {
	for {
		_, data, err := c.signalingConn.ReadMessage()
		if err != nil {
			break
		}
		c.parseSignalingMessage(data)
	}
}

func (c *ClientServer) sendSignalingMessage(in any) {
	if data, err := c.codecHelper.Encode(in); err == nil {
		_ = c.signalingConn.WriteMessage(websocket.BinaryMessage, data.ToData().B)
	}
}

func (c *ClientServer) parseSignalingMessage(data []byte) {
	if packet, err := c.codecHelper.Decode(data); err == nil {
		_ = c.dp.Call(context.Background(), nil, packet)
	}
}

func (c *ClientServer) Ping(in *proto.Ping) {
	c.sendSignalingMessage(&proto.Pong{})
}

func (c *ClientServer) Offer(in *model.Offer) {
	if client, ok := c.clients[in.UserId]; ok {
		client.Close()
		delete(c.clients, in.UserId)
	}
	if client, err := c.CreateClient(in.UserId); err == nil {
		if answer, err := client.CreateAnswer([]byte(in.Data)); err == nil {
			c.sendSignalingMessage(&model.Answer{
				UserId: in.UserId,
				Data:   string(answer),
			})
			c.clients[in.UserId] = client
		}
	}
}

func (c *ClientServer) Answer(in *model.Answer) {
	if client, ok := c.clients[in.UserId]; ok {
		_ = client.Answer([]byte(in.Data))
	}
}

func (c *ClientServer) Candidate(in *model.Candidate) {
	if client, ok := c.clients[in.UserId]; ok {
		_ = client.OnICECandidate([]byte(in.Candidate))
	}
}

func (c *ClientServer) ListRoomResponse(in *model.ListRoomResponse) {
	c.layout.UpdateRoomList(in.Rooms)
}

func (c *ClientServer) JoinRoomSuccess(in *model.JoinRoomSuccess) {
	c.layout.UpdateMessageBar("join room success", "green")
	c.layout.WriteMessage(c.layout.GetUsername(), "join room")
	c.layout.UpdateUserList(in.Users)
}

func (c *ClientServer) JoinRoomNotify(in *model.JoinRoomNotify) {
	c.layout.WriteMessage(in.Name, "join room")
	c.layout.UpdateUserList(in.Users)

	if client, ok := c.clients[in.UserId]; ok {
		client.Close()
		delete(c.clients, in.UserId)
	}

	if client, err := c.CreateClient(in.UserId); err == nil {
		if offer, err := client.CreateOffer(); err == nil {
			c.sendSignalingMessage(&model.Offer{
				UserId: in.UserId,
				Data:   string(offer),
			})
			c.clients[in.UserId] = client
		}
	}
}

func (c *ClientServer) LeaveRoomNotify(in *model.LeaveRoomNotify) {
	c.layout.WriteMessage(in.Name, "leave room")
	c.layout.UpdateUserList(in.Users)

	if peer, ok := c.clients[in.UserId]; ok {
		//delete peer connection
		peer.Close()
		delete(c.clients, in.UserId)
	}
}

func (c *ClientServer) ListRoomUsersResponse(in *model.ListRoomUsersResponse) {
}

func (c *ClientServer) CreateClient(id int64) (*p2p.P2PClient, error) {
	client := p2p.NewP2PClient(id)
	if err := client.Create(); err != nil {
		return nil, err
	}
	client.OnCandidate(func(id int64, s string) {
		c.sendSignalingMessage(&model.Candidate{
			UserId:    id,
			Candidate: s,
		})
	})
	client.OnMessage(func(b []byte) {
		var msg model.Message
		if err := gproto.Unmarshal(b, &msg); err == nil {
			c.layout.WriteMessage(msg.Name, msg.Data)
		}
	})
	client.OnClose(func(id int64) {
		client.Close()
		delete(c.clients, id)
	})

	return client, nil
}
