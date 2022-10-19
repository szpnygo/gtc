package gtcserver

import (
	"context"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/session"
	"github.com/szpnygo/gtc/server/internal/svc"
	"github.com/szpnygo/gtc/server/model"
	"go.uber.org/zap"
)

type LeaveRoomLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewLeaveRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *LeaveRoomLogic {
	return &LeaveRoomLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("leave_room_logic"),
	}
}

func (l *LeaveRoomLogic) Handler(in *model.LeaveRoom) {
	group, exist := l.svcCtx.GS.GetGroup(in.RoomId)
	if !exist {
		// the room does not exist
		return
	}

	// leave room
	if err := group.RemoveUser(l.ctx, l.session.UID()); err == nil {
		users := group.GetUsers(l.ctx)

		gogs.BroadcastMessage(users, &model.LeaveRoomNotify{
			RoomId: group.GetGroupName(l.ctx),
			UserId: l.session.ID(),
			Name:   in.Name,
			Users:  filterUid(users),
		}, nil)
	}
}
