package main

import (
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	return
}

func runClient() {
	Conn, err := grpc.Dial("localhost:50051")
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}
	defer func(Conn *grpc.ClientConn) {
		err = Conn.Close()
		if err != nil {
			zerolog.Fatal().Msg(err.Error())
			return
		}
	}(Conn)

	_ = pb.NewBlogSystemClient(Conn)

}
