// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: attendance.proto

package attendance

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

// AttendanceClient is the client API for Attendance service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AttendanceClient interface {
	Add(ctx context.Context, in *AddReq, opts ...grpc.CallOption) (*AddResp, error)
}

type attendanceClient struct {
	cc grpc.ClientConnInterface
}

func NewAttendanceClient(cc grpc.ClientConnInterface) AttendanceClient {
	return &attendanceClient{cc}
}

func (c *attendanceClient) Add(ctx context.Context, in *AddReq, opts ...grpc.CallOption) (*AddResp, error) {
	out := new(AddResp)
	err := c.cc.Invoke(ctx, "/add.attendance/add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AttendanceServer is the server API for Attendance service.
// All implementations must embed UnimplementedAttendanceServer
// for forward compatibility
type AttendanceServer interface {
	Add(context.Context, *AddReq) (*AddResp, error)
	mustEmbedUnimplementedAttendanceServer()
}

// UnimplementedAttendanceServer must be embedded to have forward compatible implementations.
type UnimplementedAttendanceServer struct {
}

func (UnimplementedAttendanceServer) Add(context.Context, *AddReq) (*AddResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedAttendanceServer) mustEmbedUnimplementedAttendanceServer() {}

// UnsafeAttendanceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AttendanceServer will
// result in compilation errors.
type UnsafeAttendanceServer interface {
	mustEmbedUnimplementedAttendanceServer()
}

func RegisterAttendanceServer(s grpc.ServiceRegistrar, srv AttendanceServer) {
	s.RegisterService(&Attendance_ServiceDesc, srv)
}

func _Attendance_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttendanceServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/add.attendance/add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttendanceServer).Add(ctx, req.(*AddReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Attendance_ServiceDesc is the grpc.ServiceDesc for Attendance service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Attendance_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "add.attendance",
	HandlerType: (*AttendanceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "add",
			Handler:    _Attendance_Add_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "attendance.proto",
}
