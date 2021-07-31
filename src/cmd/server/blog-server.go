package main

import (
	"github.com/DapperBlondie/blog-system/src/cmd/server/db"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"net"
	"os"
	"os/signal"
)

type Config struct {
	MongoDB *db.MDatabase
}

var aC *Config

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
	defer func(listener net.Listener) {
		err = listener.Close()
		if err != nil {

		}
	}(listener)

	// Use signal pkg for interrupting our server with CTRL+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	srv := grpc.NewServer()
	defer srv.Stop()

	pb.RegisterBlogSystemServer(srv, &BlogSystem{})

	aC = &Config{MongoDB: &db.MDatabase{
		MSession:     nil,
		Mdb:          nil,
		MCollections: make(map[string]*mgo.Collection),
	}}

	aC.MongoDB.MSession, err = db.NewSession("localhost:27017")
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}
	aC.MongoDB.AddDatabase("blog_system")
	aC.MongoDB.AddCollection("blogs")

	go func() {
		zerolog.Print("Blog gRPC server is listening on localhost:50051 ...")
		err = srv.Serve(listener)
		if err != nil {
			zerolog.Fatal().Msg(err.Error() + "; Occurred in serving on :50051")
			return
		}

	}()

	<-sigChan

	zerolog.Log().Msg("Blog Server was interrupted.")
	return nil
}
