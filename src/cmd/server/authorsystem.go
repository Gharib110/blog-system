package main

import (
	"context"
	"github.com/DapperBlondie/blog-system/src/cmd/server/db"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

type AuthorSystem struct {
	SignalChan  chan error
	OkChan      chan bool
	UpdateMutex *sync.RWMutex
	DeleteMutex *sync.RWMutex
}

// CreateAuthor use for creating db.Author and pb.CreateAuthorResponse
func (as *AuthorSystem) CreateAuthor(ctx context.Context, r *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
	var rspAuthor *pb.CreateAuthorResponse

	go func() {
		author := &db.Author{
			ID:     bson.NewObjectId(),
			Name:   r.GetAuthor().Name,
			Career: r.GetAuthor().Career,
		}

		err := aC.MongoDB.MCollections["authors"].Insert(author)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			as.SignalChan <- status.Error(codes.Internal, err.Error())
			return
		}

		rspAuthor = &pb.CreateAuthorResponse{
			Author: &pb.Author{
				Id:     author.ID.Hex(),
				Name:   author.Name,
				Career: author.Career,
			},
		}
		as.OkChan <- true

		return
	}()

	select {
	case <-ctx.Done():
		return nil, status.Error(status.Code(ctx.Err()), ctx.Err().Error())
	case err := <-as.SignalChan:
		return nil, err
	case <-as.OkChan:
		return rspAuthor, nil
	}
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
