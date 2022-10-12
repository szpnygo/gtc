package gtcserver

import (
	"context"

	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/session"
	"github.com/szpnygo/gtc/server/internal/message"
	"github.com/szpnygo/gtc/server/internal/svc"
	"github.com/szpnygo/gtc/server/model"
	"go.uber.org/zap"
)

type ListRoomUsersLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewListRoomUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *ListRoomUsersLogic {
	return &ListRoomUsersLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("list_room_users_logic"),
	}
}

func (l *ListRoomUsersLogic) Handler(in *model.ListRoomUsers) {
	group, exist := l.svcCtx.GS.GetGroup(in.RoomId)
	if !exist {
		// the room does not exist
		return
	}

	_ = message.SendListRoomUsersResponse(l.session, &model.ListRoomUsersResponse{
		RoomId: group.GetGroupName(l.ctx),
		Users:  filterUid(group.GetUsers(l.ctx)),
	})
}
