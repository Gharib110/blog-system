package main

import (
	"context"
	"github.com/DapperBlondie/blog-system/src/cmd/server/db"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
	"os"
	"os/signal"
	"sync"
)

type Config struct {
	MongoDB    *db.MDatabase
	SignalChan chan error
	sigMutex   *sync.Mutex
	okChan     chan bool
	okMutex    *sync.Mutex
}

var aC *Config

type BlogSystem struct {
}

func (b BlogSystem) CreateBlog(ctx context.Context, r *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	var respBlog *pb.CreateBlogResponse

	go func() {
		tBlog := r.GetBlog()
		blog := &db.BlogItem{
			ID:       bson.NewObjectId(),
			AuthorID: tBlog.AuthorId,
			Content:  tBlog.Content,
			Title:    tBlog.Title,
		}

		err := aC.MongoDB.MCollections["blogs"].Insert(blog)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			aC.sigMutex.Lock()
			aC.SignalChan <- status.Error(codes.Internal, err.Error())
			aC.sigMutex.Unlock()
		}

		respBlog = &pb.CreateBlogResponse{Blog: &pb.Blog{
			Id:       blog.ID.Hex(),
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		}}

		aC.okMutex.Lock()
		aC.okChan <- true
		aC.okMutex.Unlock()
		return
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		return nil, status.Error(status.Code(err), err.Error())
	case err := <-aC.SignalChan:
		return nil, err
	case <-aC.okChan:
		return respBlog, nil
	}
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

	aC = &Config{
		MongoDB: &db.MDatabase{
			MSession:     nil,
			Mdb:          nil,
			MCollections: make(map[string]*mgo.Collection),
		},
		SignalChan: make(chan error, 10),
		sigMutex:   &sync.Mutex{},
		okChan:     make(chan bool, 10),
		okMutex:    &sync.Mutex{},
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
