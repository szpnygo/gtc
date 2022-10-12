package gtcserver

import (
	"context"

	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/session"
	"github.com/szpnygo/gtc/server/internal/svc"
	"github.com/szpnygo/gtc/server/model"
	"go.uber.org/zap"
)

type OfferLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewOfferLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *OfferLogic {
	return &OfferLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("offer_logic"),
	}
}

func (l *OfferLogic) Handler(in *model.Offer) {
	if sess, err := session.GetSessionByID(in.UserId); err == nil {
		in.UserId = l.session.ID()
		_ = sess.SendMessage(in)
	}
}
