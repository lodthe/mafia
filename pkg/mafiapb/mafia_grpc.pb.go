// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: pkg/mafiapb/mafia.proto

package mafiapb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MafiaClient is the client API for Mafia service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MafiaClient interface {
	JoinGame(ctx context.Context, in *JoinGameRequest, opts ...grpc.CallOption) (Mafia_JoinGameClient, error)
	GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error)
	GetPlayersWithRoles(ctx context.Context, in *GetPlayersWithRolesRequest, opts ...grpc.CallOption) (*GetPlayersWithRolesResponse, error)
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	DayVote(ctx context.Context, in *DayVoteRequest, opts ...grpc.CallOption) (*DayVoteResponse, error)
	NightVote(ctx context.Context, in *NightVoteRequest, opts ...grpc.CallOption) (*NightVoteResponse, error)
	CheckTeam(ctx context.Context, in *CheckTeamRequest, opts ...grpc.CallOption) (*CheckTeamResponse, error)
}

type mafiaClient struct {
	cc grpc.ClientConnInterface
}

func NewMafiaClient(cc grpc.ClientConnInterface) MafiaClient {
	return &mafiaClient{cc}
}

func (c *mafiaClient) JoinGame(ctx context.Context, in *JoinGameRequest, opts ...grpc.CallOption) (Mafia_JoinGameClient, error) {
	stream, err := c.cc.NewStream(ctx, &Mafia_ServiceDesc.Streams[0], "/mafia.Mafia/JoinGame", opts...)
	if err != nil {
		return nil, err
	}
	x := &mafiaJoinGameClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Mafia_JoinGameClient interface {
	Recv() (*GameEvent, error)
	grpc.ClientStream
}

type mafiaJoinGameClient struct {
	grpc.ClientStream
}

func (x *mafiaJoinGameClient) Recv() (*GameEvent, error) {
	m := new(GameEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mafiaClient) GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error) {
	out := new(GetStatusResponse)
	err := c.cc.Invoke(ctx, "/mafia.Mafia/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mafiaClient) GetPlayersWithRoles(ctx context.Context, in *GetPlayersWithRolesRequest, opts ...grpc.CallOption) (*GetPlayersWithRolesResponse, error) {
	out := new(GetPlayersWithRolesResponse)
	err := c.cc.Invoke(ctx, "/mafia.Mafia/GetPlayersWithRoles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mafiaClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, "/mafia.Mafia/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mafiaClient) DayVote(ctx context.Context, in *DayVoteRequest, opts ...grpc.CallOption) (*DayVoteResponse, error) {
	out := new(DayVoteResponse)
	err := c.cc.Invoke(ctx, "/mafia.Mafia/DayVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mafiaClient) NightVote(ctx context.Context, in *NightVoteRequest, opts ...grpc.CallOption) (*NightVoteResponse, error) {
	out := new(NightVoteResponse)
	err := c.cc.Invoke(ctx, "/mafia.Mafia/NightVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mafiaClient) CheckTeam(ctx context.Context, in *CheckTeamRequest, opts ...grpc.CallOption) (*CheckTeamResponse, error) {
	out := new(CheckTeamResponse)
	err := c.cc.Invoke(ctx, "/mafia.Mafia/CheckTeam", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MafiaServer is the server API for Mafia service.
// All implementations must embed UnimplementedMafiaServer
// for forward compatibility
type MafiaServer interface {
	JoinGame(*JoinGameRequest, Mafia_JoinGameServer) error
	GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error)
	GetPlayersWithRoles(context.Context, *GetPlayersWithRolesRequest) (*GetPlayersWithRolesResponse, error)
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	DayVote(context.Context, *DayVoteRequest) (*DayVoteResponse, error)
	NightVote(context.Context, *NightVoteRequest) (*NightVoteResponse, error)
	CheckTeam(context.Context, *CheckTeamRequest) (*CheckTeamResponse, error)
	mustEmbedUnimplementedMafiaServer()
}

// UnimplementedMafiaServer must be embedded to have forward compatible implementations.
type UnimplementedMafiaServer struct {
}

func (UnimplementedMafiaServer) JoinGame(*JoinGameRequest, Mafia_JoinGameServer) error {
	return status.Errorf(codes.Unimplemented, "method JoinGame not implemented")
}
func (UnimplementedMafiaServer) GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedMafiaServer) GetPlayersWithRoles(context.Context, *GetPlayersWithRolesRequest) (*GetPlayersWithRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlayersWithRoles not implemented")
}
func (UnimplementedMafiaServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMafiaServer) DayVote(context.Context, *DayVoteRequest) (*DayVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DayVote not implemented")
}
func (UnimplementedMafiaServer) NightVote(context.Context, *NightVoteRequest) (*NightVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NightVote not implemented")
}
func (UnimplementedMafiaServer) CheckTeam(context.Context, *CheckTeamRequest) (*CheckTeamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckTeam not implemented")
}
func (UnimplementedMafiaServer) mustEmbedUnimplementedMafiaServer() {}

// UnsafeMafiaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MafiaServer will
// result in compilation errors.
type UnsafeMafiaServer interface {
	mustEmbedUnimplementedMafiaServer()
}

func RegisterMafiaServer(s grpc.ServiceRegistrar, srv MafiaServer) {
	s.RegisterService(&Mafia_ServiceDesc, srv)
}

func _Mafia_JoinGame_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(JoinGameRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MafiaServer).JoinGame(m, &mafiaJoinGameServer{stream})
}

type Mafia_JoinGameServer interface {
	Send(*GameEvent) error
	grpc.ServerStream
}

type mafiaJoinGameServer struct {
	grpc.ServerStream
}

func (x *mafiaJoinGameServer) Send(m *GameEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _Mafia_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mafia.Mafia/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServer).GetStatus(ctx, req.(*GetStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mafia_GetPlayersWithRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlayersWithRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServer).GetPlayersWithRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mafia.Mafia/GetPlayersWithRoles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServer).GetPlayersWithRoles(ctx, req.(*GetPlayersWithRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mafia_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mafia.Mafia/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mafia_DayVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServer).DayVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mafia.Mafia/DayVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServer).DayVote(ctx, req.(*DayVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mafia_NightVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NightVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServer).NightVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mafia.Mafia/NightVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServer).NightVote(ctx, req.(*NightVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mafia_CheckTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServer).CheckTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mafia.Mafia/CheckTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServer).CheckTeam(ctx, req.(*CheckTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Mafia_ServiceDesc is the grpc.ServiceDesc for Mafia service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mafia_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mafia.Mafia",
	HandlerType: (*MafiaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _Mafia_GetStatus_Handler,
		},
		{
			MethodName: "GetPlayersWithRoles",
			Handler:    _Mafia_GetPlayersWithRoles_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _Mafia_SendMessage_Handler,
		},
		{
			MethodName: "DayVote",
			Handler:    _Mafia_DayVote_Handler,
		},
		{
			MethodName: "NightVote",
			Handler:    _Mafia_NightVote_Handler,
		},
		{
			MethodName: "CheckTeam",
			Handler:    _Mafia_CheckTeam_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "JoinGame",
			Handler:       _Mafia_JoinGame_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/mafiapb/mafia.proto",
}