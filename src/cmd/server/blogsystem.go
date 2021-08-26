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

// BlogSystem Holder of RPC methods for blogSystem service
type BlogSystem struct {
	SignalChan  chan error
	OkChan      chan bool
	UpdateMutex *sync.RWMutex
	DeleteMutex *sync.RWMutex
}

// ListBlog use for getting several blogs by specify the number of blogs that we want
func (b *BlogSystem) ListBlog(r *pb.ListBlogRequest, stream pb.BlogSystem_ListBlogServer) error {
	blogI := &db.BlogItem{
		ID:       "",
		AuthorID: "",
		Content:  "",
		Title:    "",
	}
	lsBlog := &pb.ListBlogResponse{
		Blog: &pb.Blog{
			Id:       "",
			AuthorId: "",
			Title:    "",
			Content:  "",
		},
	}

	go func() {
		b.UpdateMutex.RLock()
		iterator := aC.MongoDB.MCollections["blogs"].Find(nil).Limit(int(r.GetBlogSignal())).Iter()
		b.UpdateMutex.RUnlock()
		defer func(iterator *mgo.Iter) {
			err := iterator.Close()
			if err != nil {
				zerolog.Error().Msg(err.Error() + "; Error in closing the mongodb iterator")
				return
			}
		}(iterator)

		for !iterator.Done() {
			sigB := iterator.Next(blogI)
			if sigB {
				lsBlog.GetBlog().Id = blogI.ID.Hex()
				lsBlog.GetBlog().Content = blogI.Content
				lsBlog.GetBlog().Title = blogI.Title
				lsBlog.GetBlog().AuthorId = blogI.AuthorID

				err := stream.Send(lsBlog)
				if err != nil {
					zerolog.Error().Msg(err.Error())
					b.SignalChan <- status.Error(status.Code(err), err.Error()+
						"; An Internal Error occurred in sending response to client")

					return
				}
			} else {
				b.SignalChan <- status.Error(codes.Internal, "Error in unmarshalling the data; Occurred in GetAllBlogs")
				zerolog.Error().Msg("Error in unmarshalling the data")

				return
			}
		}

		b.OkChan <- true
		return
	}()

	select {
	case <-stream.Context().Done():
		err := stream.Context().Err()
		return status.Error(status.Code(err), "; An internal error occurred in streaming cause of stream.Context ")
	case err := <-b.SignalChan:
		return err
	case <-b.OkChan:
		return nil
	}
}

// DeleteBlog use for deleting a blog by its own id
func (b *BlogSystem) DeleteBlog(ctx context.Context, r *pb.DeleteBlogRequest) (*pb.DeleteBlogResponse, error) {
	var resp *pb.DeleteBlogResponse

	go func() {
		b.DeleteMutex.Lock()
		err := aC.MongoDB.MCollections["blogs"].Remove(bson.M{"_id": bson.ObjectIdHex(r.GetBlogId())})
		b.DeleteMutex.Unlock()
		if err != nil {
			zerolog.Error().Msg(err.Error() + "; Occurred in deleting a blog with ID")
			b.SignalChan <- status.Error(codes.Internal, err.Error())
			return
		}

		b.OkChan <- true
		return
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		return nil, status.Error(status.Code(err), err.Error())
	case err := <-b.SignalChan:
		return nil, status.Error(status.Code(err), err.Error())
	case <-b.OkChan:
		resp = &pb.DeleteBlogResponse{BlogId: r.GetBlogId()}
		return resp, nil
	}
}

// UpdateBlog use for updating a blog with its own ID
func (b *BlogSystem) UpdateBlog(ctx context.Context, r *pb.UpdateBlogRequest) (*pb.UpdateBlogResponse, error) {
	var respBlog *pb.UpdateBlogResponse
	var blogItem *db.BlogItem

	go func() {
		b.UpdateMutex.RLock()
		err := aC.MongoDB.MCollections["blogs"].Find(bson.M{"_id": bson.ObjectIdHex(r.GetBlog().GetId())}).One(&blogItem)
		b.UpdateMutex.RUnlock()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			b.SignalChan <- status.Error(codes.Unavailable, err.Error())

			return
		}

		blogItem.ID = bson.ObjectIdHex(r.GetBlog().GetId())
		blogItem.Title = r.GetBlog().GetTitle()
		blogItem.Content = r.GetBlog().GetContent()
		blogItem.AuthorID = r.GetBlog().GetAuthorId()

		b.UpdateMutex.Lock()
		err = aC.MongoDB.MCollections["blogs"].UpdateId(bson.M{"_id": bson.ObjectIdHex(r.GetBlog().GetId())}, blogItem)
		b.UpdateMutex.Unlock()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			b.SignalChan <- status.Error(codes.Internal, err.Error())
			return
		}

		b.OkChan <- true
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		return nil, status.Error(status.Code(err), err.Error())
	case err := <-b.SignalChan:
		return nil, status.Error(status.Code(err), err.Error())
	case <-b.OkChan:
		respBlog = &pb.UpdateBlogResponse{Blog: &pb.Blog{
			Id:       blogItem.ID.Hex(),
			AuthorId: blogItem.AuthorID,
			Title:    blogItem.Title,
			Content:  blogItem.Content,
		}}
		return respBlog, nil
	}
}

// ReadBlog use for reading a blog with its ID
func (b *BlogSystem) ReadBlog(ctx context.Context, r *pb.ReadBlogRequest) (*pb.ReadBlogResponse, error) {
	var respBlog *pb.ReadBlogResponse
	var blogItem *db.BlogItem

	go func() {
		b.DeleteMutex.RLock()
		err := aC.MongoDB.MCollections["blogs"].Find(bson.M{"_id": bson.ObjectIdHex(r.GetBlogId())}).One(&blogItem)
		b.DeleteMutex.RUnlock()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			b.SignalChan <- status.Error(codes.Unavailable, err.Error())

			return
		}
		respBlog = &pb.ReadBlogResponse{Blog: &pb.Blog{
			Id:       blogItem.ID.Hex(),
			AuthorId: blogItem.AuthorID,
			Title:    blogItem.Title,
			Content:  blogItem.Content,
		}}

		b.OkChan <- true
		return
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		return nil, status.Error(status.Code(err), err.Error())
	case err := <-b.SignalChan:
		return nil, status.Error(status.Code(err), err.Error())
	case <-b.OkChan:
		return respBlog, nil
	}
}

// CreateBlog use for creating blog
func (b *BlogSystem) CreateBlog(ctx context.Context, r *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	var respBlog *pb.CreateBlogResponse

	go func() {
		tBlog := r.GetBlog()
		blog := &db.BlogItem{
			ID:       bson.NewObjectId(),
			AuthorID: tBlog.AuthorId,
			Content:  tBlog.Content,
			Title:    tBlog.Title,
		}

		err := aC.MongoDB.MCollections["blogs"].Insert(blog)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			b.SignalChan <- status.Error(codes.Internal, err.Error())

			return
		}

		respBlog = &pb.CreateBlogResponse{Blog: &pb.Blog{
			Id:       blog.ID.Hex(),
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		}}

		b.OkChan <- true

		return
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		return nil, status.Error(status.Code(err), err.Error())
	case err := <-b.SignalChan:
		return nil, err
	case <-b.OkChan:
		return respBlog, nil
	}
}
