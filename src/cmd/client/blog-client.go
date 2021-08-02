package main

import (
	"context"
	"github.com/DapperBlondie/blog-system/src/cmd/client/api"
	"github.com/DapperBlondie/blog-system/src/cmd/client/models"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// ClientConfig useful for holding the client configuration objects
type ClientConfig struct {
	ClientConn *grpc.ClientConn
	BlogClient pb.BlogSystemClient
	RestConfig *api.RestConf
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

	srv := http.Server{
		Addr:              "localhost:8080",
		Handler:           api.Routes(),
		ReadTimeout:       time.Second * 20,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 10,
	}

	idleChan := make(chan struct{}, 1)
	go handlingPrettyShutdown(&srv, idleChan)

	zerolog.Log().Msg("HTTP1.X server is listening on localhost:8080 ...")
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		zerolog.Error().Msg(err.Error())
	}

	<-idleChan
	zerolog.Log().Msg("HTTP1.X server shutdown successfully ... ")

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

// handlingPrettyShutdown use for shutdown the HTTP1.X server gracefully
func handlingPrettyShutdown(srv *http.Server, idleC chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	if err := srv.Shutdown(context.Background()); err != nil {
		zerolog.Error().Msg(err.Error())
	}

	close(idleC)
}

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

// Testing Creating blog
//resp, err := clientConfig.CreateBlogs()
//if err != nil {
//	return
//}
//fmt.Println(resp)

// Testing Reading blog with _id
//blog, err := clientConfig.ReadBlogs("61065b241d06c4042c5cb97f")
//if err != nil {
//	zerolog.Error().Msg(err.Error())
//	return
//}
//fmt.Println(blog)
