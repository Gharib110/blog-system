package main

import (
	"context"
	"github.com/DapperBlondie/blog-system/src/service/pb"
)

type AuthorSystem struct {
}

func (a *AuthorSystem) CreateAuthor(ctx context.Context, request *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
	panic("implement me")
}

func (a *AuthorSystem) ReadAuthor(ctx context.Context, request *pb.ReadAuthorRequest) (*pb.ReadAuthorResponse, error) {
	panic("implement me")
}

func (a *AuthorSystem) UpdateAuthor(ctx context.Context, request *pb.UpdateAuthorRequest) (*pb.UpdateAuthorResponse, error) {
	panic("implement me")
}

func (a *AuthorSystem) DeleteAuthor(ctx context.Context, request *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
	panic("implement me")
}

func (a *AuthorSystem) ListAuthor(request *pb.ListAuthorRequest, server pb.AuthorSystem_ListAuthorServer) error {
	panic("implement me")
}
