package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func main() {
	BlogClient, err := createClient()
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}

	err = createBlog(BlogClient)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}

func createClient() (pb.BlogSystemClient, error) {
	Conn, err := grpc.Dial("localhost:50051")
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return nil, status.Error(status.Code(err), err.Error())
	}
	defer func(Conn *grpc.ClientConn) {
		err = Conn.Close()
		if err != nil {
			zerolog.Fatal().Msg(err.Error())
			return
		}
	}(Conn)

	blogClient := pb.NewBlogSystemClient(Conn)

	return blogClient, nil
}

func createBlog(c pb.BlogSystemClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	blog := &pb.CreateBlogRequest{Blog: &pb.Blog{
		AuthorId: "",
		Title:    "",
		Content:  "",
	}}

	resp, err := c.CreateBlog(ctx, blog)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	fmt.Println(resp)
	return err
}
