package svc

import (
	"strings"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/group"
)

type ServiceContext struct {
	*gogs.App
	GS        *group.GroupServer
	GroupList []string
}

func NewServiceContext(app *gogs.App, rooms string) *ServiceContext {
	groupList := []string{"gtc", "gogs"}
	if len(rooms) > 0 {
		groupList = append(groupList, strings.Split(rooms, ",")...)
	}

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
