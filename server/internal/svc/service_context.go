package svc

import (
	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/group"
)

type ServiceContext struct {
	*gogs.App
	GS        *group.GroupServer
	GroupList []string
}

func NewServiceContext(app *gogs.App) *ServiceContext {
	groupList := []string{"gtc", "gtc dev", "golang", "open source", "gogs", "movie", "meta", "job", "city", "tech", "idea", "android", "iOS"}
	gs := group.NewGroupServer()
	for _, name := range groupList {
		gs.CreateMemoryGroup(name)
	}
	return &ServiceContext{
		App:       app,
		GS:        gs,
		GroupList: groupList,
	}
}
