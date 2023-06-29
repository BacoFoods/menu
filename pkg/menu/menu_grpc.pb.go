// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: proto/menu/menu.proto

package menu

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

const (
	MenuService_Create_FullMethodName      = "/menu.MenuService/Create"
	MenuService_List_FullMethodName        = "/menu.MenuService/List"
	MenuService_Get_FullMethodName         = "/menu.MenuService/Get"
	MenuService_Edit_FullMethodName        = "/menu.MenuService/Edit"
	MenuService_Delete_FullMethodName      = "/menu.MenuService/Delete"
	MenuService_AddCategory_FullMethodName = "/menu.MenuService/AddCategory"
)

// MenuServiceClient is the client API for MenuService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MenuServiceClient interface {
	Create(ctx context.Context, in *MenuCreateRequest, opts ...grpc.CallOption) (*Menu, error)
	List(ctx context.Context, in *MenuListRequest, opts ...grpc.CallOption) (*Menus, error)
	Get(ctx context.Context, in *MenuGetRequest, opts ...grpc.CallOption) (*Menu, error)
	Edit(ctx context.Context, in *MenuUpdateRequest, opts ...grpc.CallOption) (*Menu, error)
	Delete(ctx context.Context, in *MenuDeleteRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	AddCategory(ctx context.Context, in *MenuAddCategoryRequest, opts ...grpc.CallOption) (*Menu, error)
}

type menuServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMenuServiceClient(cc grpc.ClientConnInterface) MenuServiceClient {
	return &menuServiceClient{cc}
}

func (c *menuServiceClient) Create(ctx context.Context, in *MenuCreateRequest, opts ...grpc.CallOption) (*Menu, error) {
	out := new(Menu)
	err := c.cc.Invoke(ctx, MenuService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) List(ctx context.Context, in *MenuListRequest, opts ...grpc.CallOption) (*Menus, error) {
	out := new(Menus)
	err := c.cc.Invoke(ctx, MenuService_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) Get(ctx context.Context, in *MenuGetRequest, opts ...grpc.CallOption) (*Menu, error) {
	out := new(Menu)
	err := c.cc.Invoke(ctx, MenuService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) Edit(ctx context.Context, in *MenuUpdateRequest, opts ...grpc.CallOption) (*Menu, error) {
	out := new(Menu)
	err := c.cc.Invoke(ctx, MenuService_Edit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) Delete(ctx context.Context, in *MenuDeleteRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, MenuService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) AddCategory(ctx context.Context, in *MenuAddCategoryRequest, opts ...grpc.CallOption) (*Menu, error) {
	out := new(Menu)
	err := c.cc.Invoke(ctx, MenuService_AddCategory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MenuServiceServer is the server API for MenuService service.
// All implementations must embed UnimplementedMenuServiceServer
// for forward compatibility
type MenuServiceServer interface {
	Create(context.Context, *MenuCreateRequest) (*Menu, error)
	List(context.Context, *MenuListRequest) (*Menus, error)
	Get(context.Context, *MenuGetRequest) (*Menu, error)
	Edit(context.Context, *MenuUpdateRequest) (*Menu, error)
	Delete(context.Context, *MenuDeleteRequest) (*BasicResponse, error)
	AddCategory(context.Context, *MenuAddCategoryRequest) (*Menu, error)
	mustEmbedUnimplementedMenuServiceServer()
}

// UnimplementedMenuServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMenuServiceServer struct {
}

func (UnimplementedMenuServiceServer) Create(context.Context, *MenuCreateRequest) (*Menu, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedMenuServiceServer) List(context.Context, *MenuListRequest) (*Menus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedMenuServiceServer) Get(context.Context, *MenuGetRequest) (*Menu, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedMenuServiceServer) Edit(context.Context, *MenuUpdateRequest) (*Menu, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}
func (UnimplementedMenuServiceServer) Delete(context.Context, *MenuDeleteRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedMenuServiceServer) AddCategory(context.Context, *MenuAddCategoryRequest) (*Menu, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCategory not implemented")
}
func (UnimplementedMenuServiceServer) mustEmbedUnimplementedMenuServiceServer() {}

// UnsafeMenuServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MenuServiceServer will
// result in compilation errors.
type UnsafeMenuServiceServer interface {
	mustEmbedUnimplementedMenuServiceServer()
}

func RegisterMenuServiceServer(s grpc.ServiceRegistrar, srv MenuServiceServer) {
	s.RegisterService(&MenuService_ServiceDesc, srv)
}

func _MenuService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).Create(ctx, req.(*MenuCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).List(ctx, req.(*MenuListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).Get(ctx, req.(*MenuGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_Edit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).Edit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_Edit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).Edit(ctx, req.(*MenuUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).Delete(ctx, req.(*MenuDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_AddCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuAddCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).AddCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_AddCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).AddCategory(ctx, req.(*MenuAddCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MenuService_ServiceDesc is the grpc.ServiceDesc for MenuService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MenuService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "menu.MenuService",
	HandlerType: (*MenuServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _MenuService_Create_Handler,
		},
		{
			MethodName: "List",
			Handler:    _MenuService_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _MenuService_Get_Handler,
		},
		{
			MethodName: "Edit",
			Handler:    _MenuService_Edit_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _MenuService_Delete_Handler,
		},
		{
			MethodName: "AddCategory",
			Handler:    _MenuService_AddCategory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/menu/menu.proto",
}
