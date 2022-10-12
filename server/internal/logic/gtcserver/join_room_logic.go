package gtcserver

import (
	"context"
	"fmt"
	"strings"

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

	group, exist := l.svcCtx.GS.GetGroup(in.RoomId)
	if !exist {
		// the room does not exist
		return
	}

	// join room
	if err := group.AddUser(l.ctx, l.session.UID()); err == nil {
		users := group.GetUsers(l.ctx)
		usersFilted := filterUid(users)

		_ = message.SendJoinRoomSuccess(l.session, &model.JoinRoomSuccess{
			RoomId: group.GetGroupName(l.ctx),
			Users:  usersFilted,
		})

		session.BroadcastMessage(users, &model.JoinRoomNotify{
			RoomId: group.GetGroupName(l.ctx),
			Name:   in.Name,
			UserId: l.session.ID(),
			Users:  usersFilted,
		}, nil, l.session.UID())

		// when the session is closed, remove the user from the room
		l.session.OnClose(func(id int64) {
			_ = group.RemoveUser(l.ctx, l.session.UID())
			session.BroadcastMessage(users, &model.LeaveRoomNotify{
				RoomId: group.GetGroupName(l.ctx),
				Name:   in.Name,
				Users:  group.GetUsers(l.ctx),
			}, nil)
		})
	}

}

func filterUid(uids []string) []string {
	var users []string
	for _, uid := range uids {
		n := strings.Index(uid, "_")
		if n > 0 {
			users = append(users, uid[n+1:])
		}
	}
	return users
}
