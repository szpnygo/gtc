package server

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/codec"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/dispatch"
	"github.com/metagogs/gogs/global"
	"github.com/metagogs/gogs/proto"
	"github.com/pterm/pterm"
	"github.com/szpnygo/gtc/client/layout"
	"github.com/szpnygo/gtc/log"
	"github.com/szpnygo/gtc/p2p"
	"github.com/szpnygo/gtc/server/model"

	gproto "google.golang.org/protobuf/proto"
)

type ClientServer struct {
	api              string
	layout           *layout.LayoutManager
	done             chan struct{}
	readSignalingMsg chan []byte
	sendSignalingMsg chan any

	dp            *dispatch.DispatchServer
	codecHelper   *codec.CodecHelper
	signalingConn *websocket.Conn
	clients       map[int64]*p2p.P2PClient
	users         map[int64]string
	cacheUsers    []*model.User
	currentUserID int64
}

func NewClientServer(api string, layout *layout.LayoutManager) *ClientServer {
	global.GOGS_DISABLE_LOG = true
	dp := dispatch.NewDispatchServer()
	codecHelper := codec.NewCodecHelper(&config.Config{}, dp)

	return &ClientServer{
		dp:               dp,
		codecHelper:      codecHelper,
		layout:           layout,
		done:             make(chan struct{}),
		readSignalingMsg: make(chan []byte, 100),
		sendSignalingMsg: make(chan any, 100),
		clients:          make(map[int64]*p2p.P2PClient),
		users:            make(map[int64]string),
		api:              api,
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
		case data := <-c.readSignalingMsg:
			c.parseSignalingMessage(data)
		case data := <-c.sendSignalingMsg:
			c.sendSignalingMessage(data)
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
			if c.IsOK() {
				c.SendMessage(data)
			}
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
		c.layout.UpdateMessageBar(fmt.Sprintf("connect signaling server error %v ", err), "red")
		return
	}
	conn.SetCloseHandler(func(code int, text string) error {
		log.GTCLog.Infof("signaling server close %d %s", code, text)
		c.layout.UpdateMessageBar(fmt.Sprintf("websocket connection closed %d %s", code, text), "red")
		return nil
	})
	if err != nil {
		pterm.Error.Println(err)
		panic(err)
	}
	c.signalingConn = conn
	go c.readSignalingMessage()

	c.sendSignalingMsg <- &model.ListRoom{}
}

func (c *ClientServer) joinRoom(leave, join string) {
	c.layout.UpdateMessageBar("joining room "+join, "yellow")
	c.sendSignalingMsg <- &model.LeaveRoom{
		RoomId: leave,
		Name:   c.layout.GetUsername(),
	}
	for _, client := range c.clients {
		c.deleteClient(client.GetID())
	}
	c.sendSignalingMsg <- &model.JoinRoom{
		RoomId: join,
		Name:   c.layout.GetUsername(),
	}
}

func (c *ClientServer) readSignalingMessage() {
	for {
		if err := c.signalingConn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
			log.GTCLog.Error(err)
			c.layout.UpdateMessageBar(fmt.Sprintf("set read deadline error %v", err), "red")
		}
		_, data, err := c.signalingConn.ReadMessage()
		if err != nil {
			log.GTCLog.Error(err)
			c.layout.UpdateMessageBar(fmt.Sprintf("read message error %v", err), "red")
			break
		} else {
			c.readSignalingMsg <- data
		}

	}
}

func (c *ClientServer) sendSignalingMessage(in any) {
	log.GTCLog.Info("send singaling message")
	if data, err := c.codecHelper.Encode(in); err == nil {
		if err := c.signalingConn.WriteMessage(websocket.BinaryMessage, data.ToData().B); err != nil {
			log.GTCLog.Error(err)
			c.layout.UpdateMessageBar(fmt.Sprintf("send singaling message error %s", err), "red")
		}
	} else {
		log.GTCLog.Errorf("encode %s %v \n", reflect.TypeOf(in).Elem().Name(), err)
	}
}

func (c *ClientServer) parseSignalingMessage(data []byte) {
	if packet, err := c.codecHelper.Decode(data); err == nil {
		if err := c.dp.Call(context.Background(), nil, packet); err != nil {
			log.GTCLog.Error(err)
		}
	} else {
		log.GTCLog.Errorf("decode %v \n", err)
	}
}

func (c *ClientServer) updateUserList(users []*model.User) {
	c.cacheUsers = users
	list := []string{}
	for _, u := range users {
		client, ok := c.clients[u.Id]
		if ok && client.IsOK() || u.Id == c.currentUserID {
			list = append(list, fmt.Sprintf("<---> %s", u.Name))
		} else {
			list = append(list, fmt.Sprintf("<-x-> %s", u.Name))
		}
	}
	c.layout.UpdateUserList(list)
}

func (c *ClientServer) Ping(in *proto.Ping) {
	log.GTCLog.Println("Ping")
	if err := c.signalingConn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		log.GTCLog.Error(err)
		c.layout.UpdateMessageBar(fmt.Sprintf("set read deadline error %v", err), "red")
	}
	c.sendSignalingMsg <- &proto.Pong{}
}

func (c *ClientServer) Offer(in *model.Offer) {
	if client, ok := c.clients[in.UserId]; ok {
		c.deleteClient(client.GetID())
	}
	if client, err := c.CreateClient(in.UserId); err == nil {
		if answer, err := client.CreateAnswer([]byte(in.Data)); err == nil {
			c.sendSignalingMsg <- &model.Answer{
				UserId: in.UserId,
				Data:   string(answer),
			}
			c.addClient(in.UserId, client)
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
	log.GTCLog.Info("ListRoomResponse")
	c.layout.UpdateRoomList(in.Rooms)
}

func (c *ClientServer) JoinRoomSuccess(in *model.JoinRoomSuccess) {
	log.GTCLog.Info("JoinRoomSuccess")
	c.currentUserID = in.UserId
	c.layout.UpdateMessageBar(fmt.Sprintf("join room success %s", in.RoomId), "green")
	c.layout.WriteMessage(c.layout.GetUsername(), "join room")
	c.updateUserList(in.Users)
}

func (c *ClientServer) JoinRoomNotify(in *model.JoinRoomNotify) {
	c.users[in.UserId] = in.Name
	log.GTCLog.Info("JoinRoomNotify")
	c.layout.WriteMessage(in.Name, "join room")
	c.updateUserList(in.Users)

	if client, ok := c.clients[in.UserId]; ok {
		c.deleteClient(client.GetID())
	}

	if client, err := c.CreateClient(in.UserId); err == nil {
		if offer, err := client.CreateOffer(); err == nil {
			c.sendSignalingMsg <- &model.Offer{
				UserId: in.UserId,
				Data:   string(offer),
			}
			c.addClient(in.UserId, client)
		}
	}
}

func (c *ClientServer) LeaveRoomNotify(in *model.LeaveRoomNotify) {
	log.GTCLog.Info("LeaveRoomNotify")
	c.layout.WriteMessage(in.Name, "leave room")
	c.updateUserList(in.Users)

	c.deleteClient(in.UserId)
}

func (c *ClientServer) ListRoomUsersResponse(in *model.ListRoomUsersResponse) {
}

func (c *ClientServer) CreateClient(id int64) (*p2p.P2PClient, error) {
	client := p2p.NewP2PClient(id)
	if err := client.Create(); err != nil {
		return nil, err
	}
	client.OnCandidate(func(id int64, s string) {
		c.sendSignalingMsg <- &model.Candidate{
			UserId:    id,
			Candidate: s,
		}
	})
	client.OnMessage(func(b []byte) {
		var msg model.Message
		if err := gproto.Unmarshal(b, &msg); err == nil {
			c.layout.WriteMessage(msg.Name, msg.Data)
		}
	})
	client.OnClose(func(id int64) {
		delete(c.clients, id)
		c.layout.UpdateOnlineCount(len(c.clients))
	})
	client.OnChange(func() {
		c.updateUserList(c.cacheUsers)
	})

	return client, nil
}

func (c *ClientServer) addClient(id int64, client *p2p.P2PClient) {
	c.clients[id] = client
	c.layout.UpdateOnlineCount(len(c.clients))
}

func (c *ClientServer) deleteClient(id int64) {
	if peer, ok := c.clients[id]; ok {
		go peer.Close()
	}
	delete(c.clients, id)
	c.layout.UpdateOnlineCount(len(c.clients))
}
