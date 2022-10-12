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

type ListRoomLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewListRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *ListRoomLogic {
	return &ListRoomLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("list_room_logic"),
	}
}

func (l *ListRoomLogic) Handler(in *model.ListRoom) {
	_ = message.SendListRoomResponse(l.session, &model.ListRoomResponse{
		Rooms: l.svcCtx.GroupList,
	})
}
