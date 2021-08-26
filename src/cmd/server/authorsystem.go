package main

import (
	"context"
	"github.com/DapperBlondie/blog-system/src/cmd/server/db"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2"
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

// ReadAuthor use for reading a db.Author with its own bson.ObjectId
func (as *AuthorSystem) ReadAuthor(ctx context.Context, r *pb.ReadAuthorRequest) (*pb.ReadAuthorResponse, error) {
	var rspAuthor *pb.ReadAuthorResponse
	var authorItem *db.Author

	go func() {
		as.DeleteMutex.RLock()
		err := aC.MongoDB.MCollections["blogs"].Find(bson.M{"_id": bson.ObjectIdHex(r.GetAuthorId())}).One(&authorItem)
		as.DeleteMutex.RUnlock()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			as.SignalChan <- status.Error(codes.Internal, err.Error())
			return
		}
		rspAuthor = &pb.ReadAuthorResponse{Author: &pb.Author{
			Id:     r.GetAuthorId(),
			Name:   authorItem.Name,
			Career: authorItem.Career,
		}}

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

// UpdateAuthor use for updating specific db.Author and send back the *pb.UpdateAuthorResponse
func (as *AuthorSystem) UpdateAuthor(ctx context.Context, r *pb.UpdateAuthorRequest) (*pb.UpdateAuthorResponse, error) {
	var authorItem *db.Author

	go func() {
		as.UpdateMutex.RLock()
		err := aC.MongoDB.MCollections["authors"].Find(bson.M{"_id": bson.ObjectIdHex(r.GetAuthor().GetId())}).One(&authorItem)
		as.UpdateMutex.RUnlock()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			as.SignalChan <- status.Error(codes.Internal, err.Error())
			return
		}

		authorItem.Name = r.GetAuthor().GetName()
		authorItem.Career = r.GetAuthor().GetCareer()
		authorItem.ID = bson.ObjectIdHex(r.GetAuthor().GetId())

		as.UpdateMutex.Lock()
		err = aC.MongoDB.MCollections["blogs"].Update(bson.M{"_id": authorItem.ID}, &authorItem)
		as.UpdateMutex.Unlock()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			as.SignalChan <- status.Error(codes.Internal, err.Error())
			return
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
		rspAuthor := &pb.UpdateAuthorResponse{Author: &pb.Author{
			Id:     authorItem.ID.Hex(),
			Name:   authorItem.Name,
			Career: authorItem.Career,
		}}
		return rspAuthor, nil
	}
}

// DeleteAuthor use for deleting a pb.Author with its own bson.ObjectId
func (as *AuthorSystem) DeleteAuthor(ctx context.Context, r *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
	panic("implement me")
}

// ListAuthor use for sending multiple pb.Author with server streaming API
func (as *AuthorSystem) ListAuthor(r *pb.ListAuthorRequest, stream pb.AuthorSystem_ListAuthorServer) error {
	authorItem := &db.Author{
		ID:     "",
		Name:   "",
		Career: "",
	}
	rspAuthor := &pb.ListAuthorResponse{
		Author: &pb.Author{
			Id:     "",
			Name:   "",
			Career: "",
		},
	}

	go func() {
		as.DeleteMutex.RLock()
		iterator := aC.MongoDB.MCollections["authors"].Find(nil).Limit(int(r.GetAuthorSignal())).Iter()
		as.DeleteMutex.RUnlock()
		defer func(iterator *mgo.Iter) {
			err := iterator.Close()
			if err != nil {
				zerolog.Error().Msg(err.Error())
				return
			}
		}(iterator)

		for !iterator.Done() {
			sigB := iterator.Next(*authorItem)
			if sigB {
				rspAuthor.GetAuthor().Id = authorItem.ID.Hex()
				rspAuthor.GetAuthor().Name = authorItem.Name
				rspAuthor.GetAuthor().Career = authorItem.Career

				err := stream.Send(rspAuthor)
				if err != nil {
					zerolog.Error().Msg(err.Error())
					as.SignalChan <- status.Error(status.Code(err), err.Error()+
						"; An Internal Error occurred in sending response to client")
					return
				}
			} else {
				as.SignalChan <- status.Error(codes.Internal, "an error occurred in streaming data because, "+
					" unable to unmarshalling data or end of stream occurred")
			}
		}

		as.OkChan <- true
		return
	}()

	select {
	case <-stream.Context().Done():
		err := stream.Context().Err()
		zerolog.Error().Msg(err.Error())
		return status.Error(status.Code(err), err.Error())
	case err := <-as.SignalChan:
		zerolog.Error().Msg(err.Error())
	case <-as.OkChan:
		return nil
	}

	return nil
}
