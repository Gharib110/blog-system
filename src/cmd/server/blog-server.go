package main

import (
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

type BlogSystem struct {
}

func main() {
	err := runServer()
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}
}

func runServer() error {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		zerolog.Error().Msg(err.Error() + "; Occurred in listening to :50051")
		return err
	}

	// Use signal pkg for interrupting our server with CTRL+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	srv := grpc.NewServer()
	pb.RegisterBlogSystemServer(srv, &BlogSystem{})

	go func() {
		zerolog.Print("Blog gRPC server is listening on localhost:50051 ...")
		err = srv.Serve(listener)
		if err != nil {
			zerolog.Fatal().Msg(err.Error() + "; Occurred in serving on :50051")
			return
		}

	}()

	<-sigChan

	srv.Stop()
	err = listener.Close()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	zerolog.Log().Msg("Blog Server was interrupted.")
	return nil
}
