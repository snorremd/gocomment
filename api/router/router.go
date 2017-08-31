package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/snorremd/gocomment/api/model"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Router contains commenter used to create, update, get, and delete comments
type Router struct {
	Commenter model.CommentStore
}

type httpResponse struct {
	StatusCode  int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func jsonErrorResponse(w http.ResponseWriter, err *httpResponse) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(&err)
	return nil
}

func jsonResponse(w http.ResponseWriter, res interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httpErr := httpResponse{
			StatusCode:  http.StatusInternalServerError,
			Message:     http.StatusText(http.StatusInternalServerError),
			Description: "Could not serialize response.",
		}
		jsonErrorResponse(w, &httpErr)
		return
	}
}

func validateComment(r *http.Request) (*model.Comment, *httpResponse) {
	emptyComment := model.Comment{}
	comment := model.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		httpErr := httpResponse{
			Message:     http.StatusText(http.StatusBadRequest),
			StatusCode:  http.StatusBadRequest,
			Description: "Could not decode comment in payload.",
		}
		return nil, &httpErr
	}

	if reflect.DeepEqual(emptyComment, comment) {
		httpErr := httpResponse{
			Message:     http.StatusText(http.StatusBadRequest),
			StatusCode:  http.StatusBadRequest,
			Description: "Comment cannot be empty.",
		}
		return nil, &httpErr
	}

	if err := model.Validate(&comment); err != nil {
		httpErr := httpResponse{
			Message:     http.StatusText(http.StatusBadRequest),
			StatusCode:  http.StatusBadRequest,
			Description: "Comment contains illegal fields.",
		}
		return nil, &httpErr
	}

	return &comment, nil
}

func validateIDParam(r *http.Request) (*uint, *httpResponse) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		httpErr := &httpResponse{
			StatusCode:  http.StatusBadRequest,
			Message:     http.StatusText(http.StatusBadRequest),
			Description: fmt.Sprintf("Bad ID parameter %v.", vars["id"]),
		}
		return nil, httpErr
	}

	uintid := uint(id)

	return &uintid, nil
}

func (router *Router) commentHandlerPost(w http.ResponseWriter, r *http.Request) {
	comment, httpErr := validateComment(r)

	if httpErr != nil {
		jsonErrorResponse(w, httpErr)
		return
	}

	comment, err := router.Commenter.CreateComment(comment)

	if err != nil {
		httpErr := &httpResponse{
			StatusCode:  http.StatusInternalServerError,
			Message:     http.StatusText(http.StatusInternalServerError),
			Description: "Failed to create comment.",
		}
		jsonErrorResponse(w, httpErr)
		return
	}

	jsonResponse(w, comment, 200)
}

func (router *Router) commentHandlerGet(w http.ResponseWriter, r *http.Request) {
	id, httpErr := validateIDParam(r)

	if httpErr != nil {
		jsonErrorResponse(w, httpErr)
		return
	}

	comment, err := router.Commenter.GetComment(*id)

	if err != nil {
		httpErr := httpResponse{
			StatusCode:  http.StatusNotFound,
			Message:     http.StatusText(http.StatusNotFound),
			Description: fmt.Sprintf("Could not find comment with id %v.", *id),
		}
		jsonErrorResponse(w, &httpErr)
		return
	}

	jsonResponse(w, comment, http.StatusOK)
}

func (router *Router) commentHandlerGetAll(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]

	comments, err := router.Commenter.GetComments(url)

	if err != nil {
		httpErr := httpResponse{
			StatusCode:  http.StatusInternalServerError,
			Message:     http.StatusText(http.StatusInternalServerError),
			Description: fmt.Sprintf("Could not get comments for url %v.", url),
		}
		jsonErrorResponse(w, &httpErr)
		return
	}

	jsonResponse(w, comments, http.StatusOK)

}

func (router *Router) commentHandlerPut(w http.ResponseWriter, r *http.Request) {
	id, httpErr := validateIDParam(r)

	if httpErr != nil {
		jsonErrorResponse(w, httpErr)
		return
	}

	comment, httpErr := validateComment(r)

	if httpErr != nil {
		jsonErrorResponse(w, httpErr)
		return
	}

	comment.ID = id

	comment, err := router.Commenter.UpdateComment(comment)

	if err != nil && err == gorm.ErrRecordNotFound {
		httpErr := httpResponse{
			StatusCode:  http.StatusNotFound,
			Message:     http.StatusText(http.StatusNotFound),
			Description: fmt.Sprintf("Could not find comment with id %v.", *id),
		}
		jsonErrorResponse(w, &httpErr)
		return
	} else if err != nil {
		httpErr := &httpResponse{
			StatusCode:  http.StatusInternalServerError,
			Message:     http.StatusText(http.StatusInternalServerError),
			Description: "Failed to update comment.",
		}
		jsonErrorResponse(w, httpErr)
		return
	}

	jsonResponse(w, comment, 200)

}

func (router *Router) commentHandlerDelete(w http.ResponseWriter, r *http.Request) {
	id, httpErr := validateIDParam(r)

	if httpErr != nil {
		jsonErrorResponse(w, httpErr)
		return
	}

	comment := &model.Comment{
		ID: id,
	}

	comment, err := router.Commenter.DeleteComment(comment)

	if err != nil && err == gorm.ErrRecordNotFound {
		httpErr := httpResponse{
			StatusCode:  http.StatusNotFound,
			Message:     http.StatusText(http.StatusNotFound),
			Description: fmt.Sprintf("Could not find comment with id %v.", *id),
		}
		jsonErrorResponse(w, &httpErr)
		return
	} else if err != nil {
		httpErr := &httpResponse{
			StatusCode:  http.StatusInternalServerError,
			Message:     http.StatusText(http.StatusInternalServerError),
			Description: "Failed to delete comment.",
		}
		jsonErrorResponse(w, httpErr)
		return
	}

	response := httpResponse{
		StatusCode:  http.StatusOK,
		Message:     http.StatusText(http.StatusOK),
		Description: "Comment successfully deleted.",
	}
	jsonResponse(w, response, response.StatusCode)

}

// Router returns new mux router for comment routes
func (router *Router) Router() *mux.Router {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", router.commentHandlerPost).Methods("POST").Queries("url", "{url}")
	muxRouter.HandleFunc("/", router.commentHandlerGetAll).Methods("GET").Queries("url", "{url}")
	muxRouter.HandleFunc("/{id}", router.commentHandlerGet).Methods("GET")
	muxRouter.HandleFunc("/{id}", router.commentHandlerPut).Methods("PUT")
	muxRouter.HandleFunc("/{id}", router.commentHandlerDelete).Methods("DELETE")
	return muxRouter
}
