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
	err := createClient()
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}

	return
}

func createClient() error {
	Conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}
	defer func(Conn *grpc.ClientConn) {
		err = Conn.Close()
		if err != nil {
			zerolog.Fatal().Msg(err.Error())
			return
		}
	}(Conn)

	blogClient := pb.NewBlogSystemClient(Conn)

	blog := &pb.CreateBlogRequest{Blog: &pb.Blog{
		AuthorId: "",
		Title:    "Hello Johnny",
		Content:  "Hey Johnny, you are so fast and reliable !\nGood Job",
	}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := blogClient.CreateBlog(ctx, blog)
	respErr, ok := status.FromError(err)
	if !ok {
		zerolog.Error().Msg(respErr.Code().String() + respErr.Err().Error())
		return err
	}

	fmt.Println(resp)

	return nil
}
