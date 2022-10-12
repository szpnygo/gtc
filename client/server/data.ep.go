package server

import (
	"context"
	"reflect"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/proto"
	"github.com/metagogs/gogs/session"
	"github.com/szpnygo/gtc/server/model"
)

func RegisterAllComponents(s *gogs.App, srv Component) {
	registerGTCComponent(s, srv)

}

func registerGTCComponent(s *gogs.App, srv Component) {
	s.RegisterComponent(_GTCComponentDesc, srv)
}

type Component interface {
	Ping(in *proto.Ping)

	Offer(in *model.Offer)

	Answer(in *model.Answer)

	Candidate(in *model.Candidate)

	ListRoomResponse(in *model.ListRoomResponse)

	JoinRoomSuccess(in *model.JoinRoomSuccess)

	JoinRoomNotify(in *model.JoinRoomNotify)

	LeaveRoomNotify(in *model.LeaveRoomNotify)

	ListRoomUsersResponse(in *model.ListRoomUsersResponse)
}

func _GTCComponent_Ping_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).Ping(in.(*proto.Ping))
}

func _GTCComponent_Offer_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).Offer(in.(*model.Offer))
}

func _GTCComponent_Answer_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).Answer(in.(*model.Answer))
}

func _GTCComponent_Candidate_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).Candidate(in.(*model.Candidate))
}

func _GTCComponent_ListRoomResponse_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).ListRoomResponse(in.(*model.ListRoomResponse))
}

func _GTCComponent_JoinRoomSuccess_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).JoinRoomSuccess(in.(*model.JoinRoomSuccess))
}

func _GTCComponent_JoinRoomNotify_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).JoinRoomNotify(in.(*model.JoinRoomNotify))
}

func _GTCComponent_LeaveRoomNotify_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).LeaveRoomNotify(in.(*model.LeaveRoomNotify))
}

func _GTCComponent_ListRoomUsersResponse_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).ListRoomUsersResponse(in.(*model.ListRoomUsersResponse))
}

var _GTCComponentDesc = component.ComponentDesc{
	ComonentName:   "GTCComponent",
	ComponentIndex: 1, // equeal to module index
	ComponentType:  (*Component)(nil),
	Methods: []component.ComponentMethodDesc{
		{
			MethodIndex: packet.CreateAction(packet.SystemPacket, 1, 1),
			FieldType:   reflect.TypeOf(proto.Ping{}),
			Handler:     _GTCComponent_Ping_Handler,
			FiledHanler: func() interface{} {
				return new(proto.Ping)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.SystemPacket, 1, 2),
			FieldType:   reflect.TypeOf(proto.Pong{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(proto.Pong)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 1), // 0x810001 8454145
			FieldType:   reflect.TypeOf(model.Offer{}),
			Handler:     _GTCComponent_Offer_Handler,
			FiledHanler: func() interface{} {
				return new(model.Offer)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 2), // 0x810002 8454146
			FieldType:   reflect.TypeOf(model.Answer{}),
			Handler:     _GTCComponent_Answer_Handler,
			FiledHanler: func() interface{} {
				return new(model.Answer)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 3), // 0x810003 8454147
			FieldType:   reflect.TypeOf(model.Candidate{}),
			Handler:     _GTCComponent_Candidate_Handler,
			FiledHanler: func() interface{} {
				return new(model.Candidate)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 4), // 0x810004 8454148
			FieldType:   reflect.TypeOf(model.ListRoom{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(model.ListRoom)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 5), // 0x810005 8454149
			FieldType:   reflect.TypeOf(model.ListRoomResponse{}),
			Handler:     _GTCComponent_ListRoomResponse_Handler,
			FiledHanler: func() interface{} {
				return new(model.ListRoomResponse)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 6), // 0x810006 8454150
			FieldType:   reflect.TypeOf(model.JoinRoom{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(model.JoinRoom)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 7), // 0x810007 8454151
			FieldType:   reflect.TypeOf(model.JoinRoomSuccess{}),
			Handler:     _GTCComponent_JoinRoomSuccess_Handler,
			FiledHanler: func() interface{} {
				return new(model.JoinRoomSuccess)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 8), // 0x810008 8454152
			FieldType:   reflect.TypeOf(model.JoinRoomNotify{}),
			Handler:     _GTCComponent_JoinRoomNotify_Handler,
			FiledHanler: func() interface{} {
				return new(model.JoinRoomNotify)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 9), // 0x810009 8454153
			FieldType:   reflect.TypeOf(model.LeaveRoom{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(model.LeaveRoom)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 10), // 0x81000a 8454154
			FieldType:   reflect.TypeOf(model.LeaveRoomNotify{}),
			Handler:     _GTCComponent_LeaveRoomNotify_Handler,
			FiledHanler: func() interface{} {
				return new(model.LeaveRoomNotify)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 11), // 0x81000b 8454155
			FieldType:   reflect.TypeOf(model.ListRoomUsers{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(model.ListRoomUsers)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 12), // 0x81000c 8454156
			FieldType:   reflect.TypeOf(model.ListRoomUsersResponse{}),
			Handler:     _GTCComponent_ListRoomUsersResponse_Handler,
			FiledHanler: func() interface{} {
				return new(model.ListRoomUsersResponse)
			},
		},
	},
}
