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

type CandidateLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewCandidateLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *CandidateLogic {
	return &CandidateLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("candidate_logic"),
	}
}

func (l *CandidateLogic) Handler(in *model.Candidate) {
	if sess, err := gogs.GetSessionByID(in.UserId); err == nil {
		in.UserId = l.session.ID()
		_ = sess.SendMessage(in)
	}

}
