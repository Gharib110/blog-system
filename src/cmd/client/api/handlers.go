package api

import (
	"encoding/json"
	"github.com/DapperBlondie/blog-system/src/cmd/client/models"
	zerolog "github.com/rs/zerolog/log"
	"net/http"
)

// RestConf holding our rest api configurations
type RestConf struct {
	Mongo *MongoTools
}

var conf *RestConf

// NewRestConf use for creating the configuration structure for Rest-Api server
func NewRestConf(rc *RestConf) *RestConf {
	conf = &RestConf{
		Mongo: &MongoTools{
			MSession:    rc.Mongo.MSession,
			Mdb:         rc.Mongo.Mdb,
			MCollection: rc.Mongo.MCollection,
		},
	}

	return conf
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
func (rcf *RestConf) StatusHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Write(out)
	return
}

func (rcf *RestConf) InsertBlogHandler(w http.ResponseWriter, r *http.Request) {

}

func (rcf *RestConf) GetBlogHandler(w http.ResponseWriter, r *http.Request) {

}

func (rcf *RestConf) GetAuthorByIDHandler(w http.ResponseWriter, r *http.Request) {

}
