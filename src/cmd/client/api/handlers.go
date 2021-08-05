package api

import (
	"encoding/json"
	"github.com/DapperBlondie/blog-system/src/cmd/client/models"
	"github.com/DapperBlondie/blog-system/src/service/pb"
	"github.com/go-chi/chi"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"strconv"
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

// WriteResponseToUser a helper function for writing our json response to user
func WriteResponseToUser(w http.ResponseWriter, code int, resp interface{}) error {
	respType := reflect.TypeOf(resp).Elem()
	if respType.Kind() == reflect.Struct {
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
	} else if respType.Kind() == reflect.String {
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(code)
		_, err := w.Write([]byte(resp.(string)))
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return err
		}
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
	err := WriteResponseToUser(w, http.StatusOK, resp)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}

// InsertBlogHandler a rest api handler for inserting a blog
func (cc *ClientConfig) InsertBlogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Maybe method GET provided", http.StatusMethodNotAllowed)
		return
	}

	payload := &models.BlogItemPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		resp := &models.Status{Ok: http.StatusInternalServerError, Message: err.Error()}
		err = WriteResponseToUser(w, http.StatusInternalServerError, resp)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return
		}
	}

	createdBlog, err := cc.CreateBlogs(payload)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		resp := &models.Status{Ok: http.StatusInternalServerError, Message: err.Error()}
		err = WriteResponseToUser(w, http.StatusInternalServerError, resp)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return
		}
		return
	}
	payload.ID = bson.ObjectIdHex(createdBlog.Blog.Id)
	err = WriteResponseToUser(w, http.StatusOK, payload)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}
	return
}

// GetBlogHandler a rest api handler for get a blog by its own ID
func (cc *ClientConfig) GetBlogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Maybe method GET provided", http.StatusMethodNotAllowed)
		return
	}
	id := chi.URLParamFromCtx(r.Context(), "id")

	respBlog, err := cc.ReadBlogs(id)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		resp := &models.Status{Ok: http.StatusInternalServerError, Message: err.Error()}
		err = WriteResponseToUser(w, http.StatusInternalServerError, resp)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return
		}
		return
	}

	resp := &models.BlogItemPayload{
		ID:       bson.ObjectId(respBlog.GetBlog().GetId()),
		AuthorID: respBlog.GetBlog().GetAuthorId(),
		Content:  respBlog.GetBlog().GetContent(),
		Title:    respBlog.GetBlog().GetTitle(),
	}
	err = WriteResponseToUser(w, http.StatusOK, resp)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}

// GetAllBlogsHandler use for getting all blogs from gRPC
func (cc *ClientConfig) GetAllBlogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Maybe method GET provided", http.StatusMethodNotAllowed)
		return
	}
	numStr := chi.URLParamFromCtx(r.Context(), "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		errResp := &models.Status{
			Ok:      http.StatusInternalServerError,
			Message: err.Error(),
		}
		err := WriteResponseToUser(w, http.StatusInternalServerError, errResp)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return
		}
	}

	respBlogs, err := cc.GetAllBlogs(uint32(num))
	if err != nil {
		zerolog.Error().Msg(err.Error())
		errResp := &models.Status{
			Ok:      http.StatusInternalServerError,
			Message: err.Error(),
		}
		err := WriteResponseToUser(w, http.StatusInternalServerError, errResp)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return
		}
		return
	}

	blogs := &models.AllBlogsPayload{
		Blogs: []*models.BlogItemPayload{},
	}

	blog := &models.BlogItemPayload{
		ID:       "",
		AuthorID: "",
		Content:  "",
		Title:    "",
	}
	for _, respBlog := range respBlogs {
		blog.Title = respBlog.Blog.Title
		blog.Content = respBlog.Blog.Content
		blog.AuthorID = respBlog.Blog.AuthorId
		blog.ID = bson.ObjectIdHex(respBlog.Blog.Id)

		blogs.Blogs = append(blogs.Blogs, blog)
	}
	err = WriteResponseToUser(w, http.StatusOK, blogs)

	return
}
