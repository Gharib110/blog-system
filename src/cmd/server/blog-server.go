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
	"sync"
)

// Config use for holding the gRPC server configuration objects
type Config struct {
	MongoDB *db.MDatabase
}

var aC *Config

func main() {
	err := setupServer()
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}
}

// runServer a function for setting and configuring server Config and other stuff
func setupServer() error {
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

	pb.RegisterBlogSystemServer(srv, &BlogSystem{
		SignalChan:  make(chan error),
		OkChan:      make(chan bool),
		UpdateMutex: &sync.RWMutex{},
		DeleteMutex: &sync.RWMutex{},
	})
	pb.RegisterAuthorSystemServer(srv, &AuthorSystem{
		SignalChan:  make(chan error),
		OkChan:      make(chan bool),
		UpdateMutex: &sync.RWMutex{},
		DeleteMutex: &sync.RWMutex{},
	})

	aC = &Config{
		MongoDB: &db.MDatabase{
			MSession:     nil,
			Mdb:          nil,
			MCollections: make(map[string]*mgo.Collection),
		},
	}

	aC.MongoDB.MSession, err = db.NewSession("localhost:27017")
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}
	aC.MongoDB.AddDatabase("blog_system")
	aC.MongoDB.AddCollection("blogs")
	aC.MongoDB.AddCollection("authors")

	defer aC.MongoDB.MSession.Close()
	defer aC.MongoDB.Mdb.Logout()

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
