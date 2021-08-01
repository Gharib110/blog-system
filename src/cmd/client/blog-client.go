package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// ClientConfig useful for holding the client configuration objects
type ClientConfig struct {
	ClientConn *grpc.ClientConn
	BlogClient pb.BlogSystemClient
}

var clientConfig *ClientConfig

func main() {
	Conn, err := createClientConnection()
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}
	defer func(Conn *grpc.ClientConn) {
		err = Conn.Close()
		if err != nil {
			zerolog.Fatal().Msg(status.Error(codes.Internal, err.Error()).Error())
			return
		}
	}(Conn)

	clientConfig = &ClientConfig{ClientConn: Conn}
	blogClient := pb.NewBlogSystemClient(clientConfig.ClientConn)
	clientConfig.BlogClient = blogClient

	resp, err := clientConfig.CreateBlogs()
	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

// createClientConnection use for creating a clientConnection for our client
func createClientConnection() (*grpc.ClientConn, error) {
	Conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return nil, err
	}

	return Conn, nil
}

// CreateBlogs use for creating blog in out gRPC Server
func (cc *ClientConfig) CreateBlogs() (*pb.CreateBlogResponse, error) {
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
