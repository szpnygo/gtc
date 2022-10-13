package gtcserver

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/metagogs/gogs/group"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/session"
	"github.com/szpnygo/gtc/server/internal/message"
	"github.com/szpnygo/gtc/server/internal/svc"
	"github.com/szpnygo/gtc/server/model"
	"go.uber.org/zap"
)

type JoinRoomLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewJoinRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *JoinRoomLogic {
	return &JoinRoomLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("join_room_logic"),
	}
}

func (l *JoinRoomLogic) Handler(in *model.JoinRoom) {
	l.Logger.Info("join room", zap.String("room", in.RoomId), zap.String("name", in.Name))
	uid := fmt.Sprintf("%d_%s", l.session.ID(), in.Name)
	l.session.SetUID(uid)

	// when the session is closed, remove the user from the room
	l.session.SetOnCloseCallback(func(id int64) {
		for _, name := range l.svcCtx.GroupList {
			if g, ok := l.svcCtx.GS.GetGroup(name); ok {
				l.removeUserFromGroup(g, in.Name)
			}
		}
	})

	group, exist := l.svcCtx.GS.GetGroup(in.RoomId)
	if !exist {
		// the room does not exist
		return
	}

	// join room
	if err := group.AddUser(l.ctx, l.session.UID()); err != nil {
		return
	}

	users := group.GetUsers(l.ctx)
	usersFilted := filterUid(users)

	_ = message.SendJoinRoomSuccess(l.session, &model.JoinRoomSuccess{
		RoomId: group.GetGroupName(l.ctx),
		UserId: l.session.ID(),
		Users:  usersFilted,
	})

	session.BroadcastMessage(users, &model.JoinRoomNotify{
		RoomId: group.GetGroupName(l.ctx),
		Name:   in.Name,
		UserId: l.session.ID(),
		Users:  usersFilted,
	}, nil, l.session.UID())
}

func (l *JoinRoomLogic) removeUserFromGroup(g group.Group, name string) {
	if err := g.RemoveUser(l.ctx, l.session.UID()); err == nil {
		// broadcast the user left message if user is in the room
		list := g.GetUsers(l.ctx)
		session.BroadcastMessage(list, &model.LeaveRoomNotify{
			RoomId: g.GetGroupName(l.ctx),
			Name:   name,
			Users:  filterUid(list),
		}, nil)
	}
}

func filterUid(uids []string) []*model.User {
	var users []*model.User
	for _, uid := range uids {
		id, _ := strconv.ParseInt(strings.Split(uid, "_")[0], 10, 64)
		users = append(users, &model.User{
			Id:   id,
			Name: strings.Split(uid, "_")[1],
		})
	}

	return users
}
