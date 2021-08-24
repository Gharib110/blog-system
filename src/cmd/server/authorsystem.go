package main

import (
	"context"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	"sync"
)

type AuthorSystem struct{
	SignalChan  chan error
	OkChan      chan bool
	UpdateMutex *sync.RWMutex
	DeleteMutex *sync.Mutex
}

func (as *AuthorSystem) CreateAuthor(ctx context.Context, r *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
	panic("implement me")
}

func (as *AuthorSystem) ReadAuthor(ctx context.Context, r *pb.ReadAuthorRequest) (*pb.ReadAuthorResponse, error) {
	panic("implement me")
}

func (as *AuthorSystem) UpdateAuthor(ctx context.Context, r *pb.UpdateAuthorRequest) (*pb.UpdateAuthorResponse, error) {
	panic("implement me")
}

func (as *AuthorSystem) DeleteAuthor(ctx context.Context, r *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
	panic("implement me")
}

func (as *AuthorSystem) ListAuthor(request *pb.ListAuthorRequest, stream pb.AuthorSystem_ListAuthorServer) error {
	panic("implement me")
}
