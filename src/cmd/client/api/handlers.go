package api

import (
	"encoding/json"
	"github.com/DapperBlondie/blog-system/src/cmd/client/models"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/http"
)

// RestConf holding our rest api configurations
type RestConf struct {
	Mongo *MongoTools
}

// ClientConfig useful for holding the client configuration objects
type ClientConfig struct {
	ClientConn *grpc.ClientConn
	BlogClient pb.BlogSystemClient
	RestConfig *RestConf
}

var conf *ClientConfig

// NewClientConfig use for creating the configuration structure for whole RPC Client
func NewClientConfig(rc *ClientConfig) {
	conf = rc
}

// WriteToRestClient a helper function for writing our json response to rest client
func WriteToRestClient(w http.ResponseWriter, code int, resp *models.BlogItemPayload) error {
	out, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(out)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	return nil
}

// StatusHandler just use for showing the status of our API
func (cc *ClientConfig) StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, r.Method+"; We need GET request", http.StatusMethodNotAllowed)
		return
	}

	resp := &models.Status{
		Ok:      http.StatusOK,
		Message: "Everything is Alright !",
	}

	out, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(out)
	return
}

func (cc *ClientConfig) InsertBlogHandler(w http.ResponseWriter, r *http.Request) {

}

func (cc *ClientConfig) GetBlogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Maybe method GET provided", http.StatusMethodNotAllowed)
		return
	}
	//id := chi.URLParamFromCtx(r.Context(), "id")

}

func (cc *ClientConfig) GetAuthorByIDHandler(w http.ResponseWriter, r *http.Request) {

}

func (cc *ClientConfig) InsertAuthorHandler(w http.ResponseWriter, r *http.Request) {

}

func (cc *ClientConfig) GetAllBlogsHandler(w http.ResponseWriter, r *http.Request) {

}
