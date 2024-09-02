// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: proto/shortener.proto

package proto

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

// URLShortenerClient is the client API for URLShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerClient interface {
	SaveURL(ctx context.Context, in *SaveURLRequest, opts ...grpc.CallOption) (*SaveURLResponse, error)
	SaveJSON(ctx context.Context, in *SaveJSONRequest, opts ...grpc.CallOption) (*SaveJSONResponse, error)
	SaveJSONBatch(ctx context.Context, in *SaveJSONBatchRequest, opts ...grpc.CallOption) (*SaveJSONBatchResponse, error)
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	GetAllURLs(ctx context.Context, in *GetAllURLsRequest, opts ...grpc.CallOption) (*GetAllURLsResponse, error)
	GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error)
	DeleteURL(ctx context.Context, in *DeleteURLRequest, opts ...grpc.CallOption) (*DeleteURLResponse, error)
}

type uRLShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerClient(cc grpc.ClientConnInterface) URLShortenerClient {
	return &uRLShortenerClient{cc}
}

func (c *uRLShortenerClient) SaveURL(ctx context.Context, in *SaveURLRequest, opts ...grpc.CallOption) (*SaveURLResponse, error) {
	out := new(SaveURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.URLShortener/SaveURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) SaveJSON(ctx context.Context, in *SaveJSONRequest, opts ...grpc.CallOption) (*SaveJSONResponse, error) {
	out := new(SaveJSONResponse)
	err := c.cc.Invoke(ctx, "/shortener.URLShortener/SaveJSON", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) SaveJSONBatch(ctx context.Context, in *SaveJSONBatchRequest, opts ...grpc.CallOption) (*SaveJSONBatchResponse, error) {
	out := new(SaveJSONBatchResponse)
	err := c.cc.Invoke(ctx, "/shortener.URLShortener/SaveJSONBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.URLShortener/GetURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetAllURLs(ctx context.Context, in *GetAllURLsRequest, opts ...grpc.CallOption) (*GetAllURLsResponse, error) {
	out := new(GetAllURLsResponse)
	err := c.cc.Invoke(ctx, "/shortener.URLShortener/GetAllURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error) {
	out := new(GetStatsResponse)
	err := c.cc.Invoke(ctx, "/shortener.URLShortener/GetStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) DeleteURL(ctx context.Context, in *DeleteURLRequest, opts ...grpc.CallOption) (*DeleteURLResponse, error) {
	out := new(DeleteURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.URLShortener/DeleteURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServer is the server API for URLShortener service.
// All implementations must embed UnimplementedURLShortenerServer
// for forward compatibility
type URLShortenerServer interface {
	SaveURL(context.Context, *SaveURLRequest) (*SaveURLResponse, error)
	SaveJSON(context.Context, *SaveJSONRequest) (*SaveJSONResponse, error)
	SaveJSONBatch(context.Context, *SaveJSONBatchRequest) (*SaveJSONBatchResponse, error)
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	GetAllURLs(context.Context, *GetAllURLsRequest) (*GetAllURLsResponse, error)
	GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error)
	DeleteURL(context.Context, *DeleteURLRequest) (*DeleteURLResponse, error)
	mustEmbedUnimplementedURLShortenerServer()
}

// UnimplementedURLShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedURLShortenerServer struct {
}

func (UnimplementedURLShortenerServer) SaveURL(context.Context, *SaveURLRequest) (*SaveURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveURL not implemented")
}
func (UnimplementedURLShortenerServer) SaveJSON(context.Context, *SaveJSONRequest) (*SaveJSONResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveJSON not implemented")
}
func (UnimplementedURLShortenerServer) SaveJSONBatch(context.Context, *SaveJSONBatchRequest) (*SaveJSONBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveJSONBatch not implemented")
}
func (UnimplementedURLShortenerServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}
func (UnimplementedURLShortenerServer) GetAllURLs(context.Context, *GetAllURLsRequest) (*GetAllURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllURLs not implemented")
}
func (UnimplementedURLShortenerServer) GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedURLShortenerServer) DeleteURL(context.Context, *DeleteURLRequest) (*DeleteURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURL not implemented")
}
func (UnimplementedURLShortenerServer) mustEmbedUnimplementedURLShortenerServer() {}

// UnsafeURLShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServer will
// result in compilation errors.
type UnsafeURLShortenerServer interface {
	mustEmbedUnimplementedURLShortenerServer()
}

func RegisterURLShortenerServer(s grpc.ServiceRegistrar, srv URLShortenerServer) {
	s.RegisterService(&URLShortener_ServiceDesc, srv)
}

func _URLShortener_SaveURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).SaveURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.URLShortener/SaveURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).SaveURL(ctx, req.(*SaveURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_SaveJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveJSONRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).SaveJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.URLShortener/SaveJSON",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).SaveJSON(ctx, req.(*SaveJSONRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_SaveJSONBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveJSONBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).SaveJSONBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.URLShortener/SaveJSONBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).SaveJSONBatch(ctx, req.(*SaveJSONBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.URLShortener/GetURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetAllURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetAllURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.URLShortener/GetAllURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetAllURLs(ctx, req.(*GetAllURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.URLShortener/GetStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetStats(ctx, req.(*GetStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_DeleteURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).DeleteURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.URLShortener/DeleteURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).DeleteURL(ctx, req.(*DeleteURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortener_ServiceDesc is the grpc.ServiceDesc for URLShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.URLShortener",
	HandlerType: (*URLShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveURL",
			Handler:    _URLShortener_SaveURL_Handler,
		},
		{
			MethodName: "SaveJSON",
			Handler:    _URLShortener_SaveJSON_Handler,
		},
		{
			MethodName: "SaveJSONBatch",
			Handler:    _URLShortener_SaveJSONBatch_Handler,
		},
		{
			MethodName: "GetURL",
			Handler:    _URLShortener_GetURL_Handler,
		},
		{
			MethodName: "GetAllURLs",
			Handler:    _URLShortener_GetAllURLs_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _URLShortener_GetStats_Handler,
		},
		{
			MethodName: "DeleteURL",
			Handler:    _URLShortener_DeleteURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}