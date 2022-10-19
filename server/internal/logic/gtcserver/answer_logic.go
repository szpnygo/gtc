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

type AnswerLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *AnswerLogic {
	return &AnswerLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("answer_logic"),
	}
}

func (l *AnswerLogic) Handler(in *model.Answer) {
	if sess, err := gogs.GetSessionByID(in.UserId); err == nil {
		in.UserId = l.session.ID()
		_ = sess.SendMessage(in)
	}
}
