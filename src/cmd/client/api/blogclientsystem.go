package api

import (
	"context"
	"github.com/DapperBlondie/blog-system/src/cmd/client/models"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

// CreateBlogs use for creating blog in out gRPC Server
func (cc *ClientConfig) CreateBlogs(bp *models.BlogItemPayload) (*pb.CreateBlogResponse, error) {
	blog := &pb.CreateBlogRequest{Blog: &pb.Blog{
		AuthorId: "",
		Title:    "Hello Gopher",
		Content:  "Hey Gopher, Golang is so fantastic !\nGood Job",
	}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := cc.BlogClient.CreateBlog(ctx, blog)
	respErr, ok := status.FromError(err)
	if !ok {
		zerolog.Error().Msg(respErr.Code().String() + respErr.Err().Error())
		return nil, err
	}

	return resp, nil
}

// ReadBlogs use for reading blogs by their own IDs
func (cc *ClientConfig) ReadBlogs(id string) (*pb.ReadBlogResponse, error) {
	req := &pb.ReadBlogRequest{BlogId: id}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resBlog, err := cc.BlogClient.ReadBlog(ctx, req)
	s, ok := status.FromError(err)
	if !ok {
		zerolog.Error().Msg(s.Code().String() + " : " + s.Message())
		return nil, err
	}

	return resBlog, nil
}

// UpdateBlogs use for updating blogs by getting blog payload
func (cc *ClientConfig) UpdateBlogs(bp *models.BlogItemPayload) (*pb.UpdateBlogResponse, error) {
	req := &pb.UpdateBlogRequest{Blog: &pb.Blog{
		Id:       bp.ID.Hex(),
		AuthorId: bp.AuthorID,
		Title:    bp.Title,
		Content:  bp.Content,
	}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	updatedBlog, err := cc.BlogClient.UpdateBlog(ctx, req)
	s, ok := status.FromError(err)
	if !ok {
		zerolog.Error().Msg(s.Code().String() + "; " + s.Message())
		return nil, err
	}

	return updatedBlog, nil
}

// DeleteBlogs use for deleting blogs by getting its own ID
func (cc *ClientConfig) DeleteBlogs(id string) (*pb.DeleteBlogResponse, error) {
	req := &pb.DeleteBlogRequest{BlogId: id}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	deletedBlog, err := cc.BlogClient.DeleteBlog(ctx, req)
	s, ok := status.FromError(err)
	if !ok {
		zerolog.Error().Msg(s.Code().String() + "; " + s.Message())
		return nil, err
	}

	return deletedBlog, nil
}

// GetAllBlogs use for getting all blogs from server
func (cc *ClientConfig) GetAllBlogs(num uint32) ([]*pb.ListBlogResponse, error) {
	req := &pb.ListBlogRequest{BlogSignal: num}
	lstBlogs := []*pb.ListBlogResponse{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*10)
	defer cancel()

	stream, err := cc.BlogClient.ListBlog(ctx, req)
	s, ok := status.FromError(err)
	if !ok {
		zerolog.Error().Msg(s.Code().String() + " " + s.Message() + "; in GetAllBlogs")
		return nil, err
	}

	for {
		blog, err := stream.Recv()
		if err == io.EOF {
			zerolog.Error().Msg(err.Error())
			return lstBlogs, err
		} else if err != nil {
			zerolog.Error().Msg(err.Error())
			return nil, err
		}
		lstBlogs = append(lstBlogs, blog)
	}

}
