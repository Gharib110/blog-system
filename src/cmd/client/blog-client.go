package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/blog-system/src/cmd/client/api"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	MONGO_HOST = "localhost:27017"
	GRPC_HOST  = "localhost:50051"
	REST_HOST  = "localhost:8080"
)

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

	clientConfig := &api.ClientConfig{ClientConn: Conn}
	blogClient := pb.NewBlogSystemClient(clientConfig.ClientConn)
	clientConfig.BlogClient = blogClient

	session, err := api.NewMSession(MONGO_HOST)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}
	restConfig := &api.RestConf{
		Mongo: &api.MongoTools{
			MSession:    session,
			MCollection: make(map[string]*mgo.Collection),
		},
	}
	restConfig.Mongo.NewMDatabase("blog-system")
	restConfig.Mongo.NewMCollection("blogs")
	restConfig.Mongo.NewMCollection("authors")
	clientConfig.RestConfig = restConfig

	api.NewClientConfig(clientConfig)

	srv := http.Server{
		Addr:              REST_HOST,
		Handler:           api.Routes(),
		ReadTimeout:       time.Second * 20,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 10,
	}

	idleChan := make(chan struct{}, 1)
	go handlingPrettyShutdown(&srv, idleChan)

	zerolog.Log().Msg(fmt.Sprintf("HTTP1.X server is listening on %s ...\n", REST_HOST))
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		zerolog.Error().Msg(err.Error())
		return
	}

	<-idleChan
	zerolog.Log().Msg("HTTP1.X server shutdown successfully ... ")

	return
}

// createClientConnection use for creating a clientConnection for our client
func createClientConnection() (*grpc.ClientConn, error) {
	Conn, err := grpc.Dial(GRPC_HOST, grpc.WithInsecure())
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
