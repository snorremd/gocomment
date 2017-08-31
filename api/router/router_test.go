package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/snorremd/gocomment/api/model"

	"github.com/jinzhu/gorm"
)

type mockCommentStore struct{}

func (c mockCommentStore) GetComments(url string) ([]*model.Comment, error) {

	id1 := uint(1)
	id2 := uint(2)
	id3 := uint(3)
	createdAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	if url == "http://example.com/posts/1" {
		return []*model.Comment{
			&model.Comment{
				ID:        &id1,
				CreatedAt: &createdAt,
				UpdatedAt: &updatedAt,
				DeletedAt: nil,
				ParentID:  0,
				Content:   "Some content",
				Upvotes:   0,
				Downvotes: 0,
				URL:       "http://example.com/posts/1",
			},
			&model.Comment{
				ID:        &id2,
				CreatedAt: &createdAt,
				UpdatedAt: &updatedAt,
				DeletedAt: nil,
				ParentID:  0,
				Content:   "Some content",
				Upvotes:   0,
				Downvotes: 0,
				URL:       "http://example.com/posts/1",
			},
			&model.Comment{
				ID:        &id3,
				CreatedAt: &createdAt,
				UpdatedAt: &updatedAt,
				DeletedAt: nil,
				ParentID:  0,
				Content:   "Some content",
				Upvotes:   0,
				Downvotes: 0,
				URL:       "http://example.com/posts/1",
			},
		}, nil
	} else if url == "http://example.com/posts/1" {
		return []*model.Comment{}, nil
	}

	return nil, errors.New("some error")

}

func (c mockCommentStore) GetComment(id uint) (*model.Comment, error) {
	if id == uint(1) {

		someID := uint(1)
		createdAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

		return &model.Comment{
			ID:        &someID,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
			DeletedAt: nil,
			ParentID:  0,
			Content:   "Some content",
			Upvotes:   0,
			Downvotes: 0,
			URL:       "http://example.com/posts/1",
		}, nil

	}

	return nil, errors.New("record not found")

}

func (c mockCommentStore) CreateComment(comment *model.Comment) (*model.Comment, error) {
	if comment.Content == "Error" {
		return nil, errors.New("Failed to create comment")
	}
	id := uint(1)
	createdAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	comment.ID = &id
	comment.CreatedAt = &createdAt
	comment.UpdatedAt = &updatedAt
	comment.DeletedAt = nil
	return comment, nil
}

func (c mockCommentStore) UpdateComment(comment *model.Comment) (*model.Comment, error) {
	if *comment.ID == uint(1) {
		originalComment, err := c.GetComment(*comment.ID)
		if err != nil {
			return nil, err
		}

		updatedAt := time.Date(1970, time.January, 2, 0, 0, 0, 0, time.UTC)

		originalComment.UpdatedAt = &updatedAt

		return originalComment, nil
	}

	return nil, gorm.ErrRecordNotFound

}

func (c mockCommentStore) DeleteComment(comment *model.Comment) (*model.Comment, error) {
	if *comment.ID == uint(1) {

		someID := uint(1)
		createdAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
		deletedAt := time.Date(1970, time.January, 2, 0, 0, 0, 0, time.UTC)

		return &model.Comment{
			ID:        &someID,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
			DeletedAt: &deletedAt,
			ParentID:  0,
			Content:   "Some content",
			Upvotes:   0,
			Downvotes: 0,
			URL:       "http://example.com/posts/1",
		}, nil

	}

	return nil, gorm.ErrRecordNotFound
}

func Test_validateComment(t *testing.T) {

	comment := &model.Comment{
		Content:  "",
		ParentID: 1,
		URL:      "Foo",
	}
	payload, _ := json.Marshal(comment)
	validRequest, _ := http.NewRequest("POST", "/?url=http://example.com", bytes.NewBuffer(payload))

	id := uint(0)
	illegalPropsComment := &model.Comment{
		ID:      &id,
		Content: "",
	}
	illegalPayload, _ := json.Marshal(illegalPropsComment)
	illegalPropsRequest, _ := http.NewRequest("POST", "/?url=http://example.com", bytes.NewBuffer(illegalPayload))

	badFormatPayload := bytes.NewBufferString("Not json")
	badlyFormattedJSONRequest, _ := http.NewRequest("POST", "/?url=http://example.com", badFormatPayload)

	tests := []struct {
		name      string
		r         *http.Request
		comment   *model.Comment
		httpError *httpResponse
	}{
		{
			name:      "Valid comment returns comment",
			r:         validRequest,
			comment:   comment,
			httpError: nil,
		},
		{
			name: "Comment with illegal props returns error",
			r:    illegalPropsRequest,
			httpError: &httpResponse{
				StatusCode:  http.StatusBadRequest,
				Message:     http.StatusText(http.StatusBadRequest),
				Description: "Comment contains illegal fields.",
			},
			comment: nil,
		},
		{
			name: "Badly formatted json payload",
			r:    badlyFormattedJSONRequest,
			httpError: &httpResponse{
				StatusCode:  http.StatusBadRequest,
				Message:     http.StatusText(http.StatusBadRequest),
				Description: "Could not decode comment in payload.",
			},
			comment: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			comment, err := validateComment(tt.r)
			if !reflect.DeepEqual(comment, tt.comment) {
				t.Errorf("validateComment() got = %v, want %v", comment, tt.comment)
			}
			if !reflect.DeepEqual(err, tt.httpError) {
				t.Errorf("validateComment() got1 = %v, want %v", err, tt.httpError)
			}
		})
	}
}

func Test_jsonResponse(t *testing.T) {
	tests := []struct {
		name            string
		res             interface{}
		statusCode      int
		contentType     string
		httpResponse    *httpResponse
		commentResponse *model.Comment
	}{
		{
			name: "Valid json encodable response",
			res: &httpResponse{
				StatusCode:  http.StatusOK,
				Message:     http.StatusText(http.StatusOK),
				Description: "Some response description.",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json; charset=utf-8",
			httpResponse: &httpResponse{
				StatusCode:  http.StatusOK,
				Message:     http.StatusText(http.StatusOK),
				Description: "Some response description.",
			},
		},
		{
			name:        "Invalid json unencodable response",
			res:         make(chan int),
			statusCode:  http.StatusOK,
			contentType: "application/json; charset=utf-8",
			httpResponse: &httpResponse{
				StatusCode:  http.StatusInternalServerError,
				Message:     http.StatusText(http.StatusInternalServerError),
				Description: "Could not serialize response.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			jsonResponse(recorder, tt.res, tt.statusCode)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected handler to respond with code %v, but got %v", tt.statusCode, recorder.Code)
			}

			if contentType := recorder.Header().Get("content-type"); contentType != tt.contentType {
				t.Errorf("Expected handler to respond with content-type %v, but got %v", tt.contentType, contentType)
			}

			if tt.commentResponse != nil { // Expect regular body
				comment := &model.Comment{}
				if err := json.NewDecoder(recorder.Body).Decode(comment); err != nil {
					t.Errorf("Could not decode comment body %v because of error %v", recorder.Body, err)
				}

				if *comment.ID != *tt.commentResponse.ID {
					t.Errorf("Expected ID %v, but was %v", tt.commentResponse.ID, comment.ID)
				}
				if comment.URL != tt.commentResponse.URL {
					t.Errorf("Expected URL %v, but was %v", tt.commentResponse.URL, comment.URL)
				}

			} else if tt.httpResponse != nil { // Expect error body
				httpResponse := &httpResponse{}
				if err := json.NewDecoder(recorder.Body).Decode(httpResponse); err != nil {
					t.Errorf("Could not decode httpError body %v because of error %v", recorder.Body, err)
				}

				if !reflect.DeepEqual(httpResponse, tt.httpResponse) {
					t.Errorf("Expected json error to be %v, but got %v", tt.httpResponse, httpResponse)
				}
			}

		})
	}
}

func Test_server_commentHandlerPost(t *testing.T) {
	router := &Router{
		Commenter: &mockCommentStore{},
	}

	id := uint(0)

	inputComment := &model.Comment{
		Content: "Some content",
	}

	outputComment, err := router.Commenter.CreateComment(&model.Comment{
		Content:   "Some content",
		Upvotes:   0,
		Downvotes: 0,
	})

	if err != nil {
		t.Errorf("Could not create output comment %v, got error %v", outputComment, err)
	}

	tests := []struct {
		name        string
		comment     *model.Comment
		url         string
		statusCode  int
		commentBody *model.Comment
		errorBody   *httpResponse
	}{
		{
			name:        "Post valid comment",
			comment:     inputComment,
			url:         "http://example.com/posts/1",
			statusCode:  http.StatusOK,
			commentBody: outputComment,
		},
		{
			name:       "Post invalid (empty) comment",
			comment:    &model.Comment{},
			url:        "http://example.com/posts/1",
			statusCode: http.StatusBadRequest,
			errorBody: &httpResponse{
				StatusCode:  http.StatusBadRequest,
				Message:     http.StatusText(http.StatusBadRequest),
				Description: "Comment cannot be empty.",
			},
		},
		{
			name: "Post invalid comment with database generated property",
			comment: &model.Comment{
				ID:      &id,
				Content: "Some content",
			},
			url:        "http://example.com/posts/1",
			statusCode: http.StatusBadRequest,
			errorBody: &httpResponse{
				StatusCode:  http.StatusBadRequest,
				Message:     http.StatusText(http.StatusBadRequest),
				Description: "Comment contains illegal fields.",
			},
		},
		{
			name: "Post comment failing database creation",
			comment: &model.Comment{
				Content: "Error",
			},
			url:        "http://example.com/posts/1",
			statusCode: http.StatusInternalServerError,
			errorBody: &httpResponse{
				StatusCode:  http.StatusInternalServerError,
				Message:     http.StatusText(http.StatusInternalServerError),
				Description: "Failed to create comment.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			payload, _ := json.Marshal(tt.comment)
			request, _ := http.NewRequest("POST", "/?url="+tt.url, bytes.NewBuffer(payload))
			recorder := httptest.NewRecorder()
			muxRouter := router.Router()
			muxRouter.ServeHTTP(recorder, request)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected handler to respond with code %v, but got %v", tt.statusCode, recorder.Code)
			}

			if tt.commentBody != nil { // Expect regular body
				comment := &model.Comment{}
				if err := json.NewDecoder(recorder.Body).Decode(comment); err != nil {
					t.Errorf("Could not decode comment body %v because of error %v", recorder.Body, err)
				}

				if *comment.ID != *tt.commentBody.ID {
					t.Errorf("Expected ID %v, but was %v", tt.commentBody.ID, comment.ID)
				}
				if comment.URL != tt.commentBody.URL {
					t.Errorf("Expected URL %v, but was %v", tt.commentBody.URL, comment.URL)
				}

			} else if tt.errorBody != nil { // Expect error body
				httpError := &httpResponse{}
				if err := json.NewDecoder(recorder.Body).Decode(httpError); err != nil {
					t.Errorf("Could not decode httpError body %v because of error %v", recorder.Body, err)
				}

				if !reflect.DeepEqual(httpError, tt.errorBody) {
					t.Errorf("Expected json error to be %v, but got %v", tt.errorBody, httpError)
				}
			}
		})
	}
}

func Test_server_commentHandlerGet(t *testing.T) {

	router := &Router{
		Commenter: &mockCommentStore{},
	}

	comment, err := router.Commenter.GetComment(uint(1))

	if err != nil {
		t.Errorf("Could not get mocked comment %v, got error %v", comment, err)
	}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name        string
		id          string
		statusCode  int
		commentBody *model.Comment
		errorBody   *httpResponse
	}{
		{
			name:        "Get comment that exists in DB",
			id:          "1",
			statusCode:  200,
			commentBody: comment,
			errorBody:   nil,
		},
		{
			name:        "Get non existing comment should cause Not Found",
			id:          "1000",
			statusCode:  404,
			commentBody: nil,
			errorBody: &httpResponse{
				StatusCode:  404,
				Message:     http.StatusText(http.StatusNotFound),
				Description: "Could not find comment with id 1000.",
			},
		},
		{
			name:        "Get comment with invalid id parameter should cause Bad Request",
			id:          "-1",
			statusCode:  400,
			commentBody: nil,
			errorBody: &httpResponse{
				StatusCode:  400,
				Message:     http.StatusText(http.StatusBadRequest),
				Description: "Bad ID parameter -1.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request, _ := http.NewRequest("GET", "/"+tt.id, nil)
			recorder := httptest.NewRecorder()
			muxRouter := router.Router()
			muxRouter.ServeHTTP(recorder, request)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected handler to respond with code %v, but got %v", tt.statusCode, recorder.Code)
			}

			if tt.commentBody != nil { // Expect regular body
				comment := &model.Comment{}
				if err := json.NewDecoder(recorder.Body).Decode(comment); err != nil {
					t.Errorf("Could not decode comment body %v because of error %v", recorder.Body, err)
				}

				if *comment.ID != *tt.commentBody.ID {
					t.Errorf("Expected ID %v, but was %v", tt.commentBody.ID, comment.ID)
				}

			} else if tt.errorBody != nil { // Expect error body
				httpError := &httpResponse{}
				if err := json.NewDecoder(recorder.Body).Decode(httpError); err != nil {
					t.Errorf("Could not decode httpError body %v because of error %v", recorder.Body, err)
				}

				if !reflect.DeepEqual(httpError, tt.errorBody) {
					t.Errorf("Expected json error to be %v, but got %v", tt.errorBody, httpError)
				}
			}
		})
	}
}

func Test_server_commentHandlerGetAll(t *testing.T) {

	router := &Router{
		Commenter: &mockCommentStore{},
	}

	comments, _ := router.Commenter.GetComments("http://example.com/posts/1")

	tests := []struct {
		name       string
		url        string
		statusCode int
		comments   []*model.Comment
		errorBody  *httpResponse
	}{
		{
			name:       "Get comments for url existing in db",
			url:        "http://example.com/posts/1",
			statusCode: http.StatusOK,
			comments:   comments,
			errorBody:  nil,
		},
		{
			name:       "Get comments when commenter returns error",
			url:        "not-in-database",
			statusCode: http.StatusInternalServerError,
			errorBody: &httpResponse{
				StatusCode:  http.StatusInternalServerError,
				Message:     http.StatusText(http.StatusInternalServerError),
				Description: "Could not get comments for url not-in-database.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request, _ := http.NewRequest("GET", "/?url="+tt.url, nil)
			recorder := httptest.NewRecorder()
			muxRouter := router.Router()
			muxRouter.ServeHTTP(recorder, request)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected handler to respond with code %v, but got %v", tt.statusCode, recorder.Code)
			}

			if tt.errorBody == nil { // Expect regular body
				comments := make([]model.Comment, 0)
				if err := json.NewDecoder(recorder.Body).Decode(&comments); err != nil {
					t.Errorf("Could not decode comments %v because of error %v", recorder.Body, err)
				}

				if reflect.DeepEqual(comments, tt.comments) {
					t.Errorf("Expected comments %v, but was %v", tt.comments, comments)
				}

			} else if tt.errorBody != nil { // Expect error body
				httpError := &httpResponse{}
				if err := json.NewDecoder(recorder.Body).Decode(httpError); err != nil {
					t.Errorf("Could not decode httpError body %v because of error %v", recorder.Body, err)
				}

				if !reflect.DeepEqual(httpError, tt.errorBody) {
					t.Errorf("Expected json error to be %v, but got %v", tt.errorBody, httpError)
				}
			}
		})
	}
}

func Test_server_commentHandlerPut(t *testing.T) {
	router := &Router{
		Commenter: &mockCommentStore{},
	}

	inputComment := &model.Comment{
		Content:   "Some content",
		Upvotes:   0,
		Downvotes: 0,
		URL:       "http://example.com/posts/1",
	}

	updatedComment, err := router.Commenter.CreateComment(&model.Comment{
		Content:   "Some content",
		Upvotes:   0,
		Downvotes: 0,
		URL:       "https://example.com/post/1",
	})

	if err != nil {
		t.Errorf("Could not create comment %v, got error %v", inputComment, err)
	}

	updatedComment, err = router.Commenter.UpdateComment(updatedComment)

	if err != nil {
		t.Errorf("Could not update comment %v, got error %v", updatedComment, err)
	}

	tests := []struct {
		name        string
		id          uint
		comment     *model.Comment
		statusCode  int
		commentBody *model.Comment
		errorBody   *httpResponse
	}{
		{
			name:        "Put valid comment",
			id:          uint(1),
			comment:     inputComment,
			statusCode:  200,
			commentBody: updatedComment,
			errorBody:   nil,
		},
		{
			name:       "Put comment no in db",
			id:         uint(1000),
			comment:    inputComment,
			statusCode: http.StatusNotFound,
			errorBody: &httpResponse{
				StatusCode:  http.StatusNotFound,
				Message:     http.StatusText(http.StatusNotFound),
				Description: "Could not find comment with id 1000.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			payload, _ := json.Marshal(tt.comment)
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/%v", tt.id), bytes.NewBuffer(payload))
			recorder := httptest.NewRecorder()
			muxRouter := router.Router()
			muxRouter.ServeHTTP(recorder, request)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected handler to respond with code %v, but got %v", tt.statusCode, recorder.Code)
			}

			if tt.commentBody != nil { // Expect regular body
				comment := &model.Comment{}
				if err := json.NewDecoder(recorder.Body).Decode(comment); err != nil {
					t.Errorf("Could not decode comment body %v because of error %v", recorder.Body, err)
				}

				if *comment.ID != *tt.commentBody.ID {
					t.Errorf("Expected ID %v, but was %v", tt.commentBody.ID, comment.ID)
				}

				if !comment.UpdatedAt.After(*tt.commentBody.CreatedAt) {
					t.Errorf("Expected UpdatedAt %v, but was %v", tt.commentBody.UpdatedAt, comment.UpdatedAt)
				}

			} else if tt.errorBody != nil { // Expect error body
				httpError := &httpResponse{}
				if err := json.NewDecoder(recorder.Body).Decode(httpError); err != nil {
					t.Errorf("Could not decode httpError body %v because of error %v", recorder.Body, err)
				}
				if !reflect.DeepEqual(httpError, tt.errorBody) {
					t.Errorf("Expected json error to be %v, but got %v", tt.errorBody, httpError)
				}
			}
		})
	}
}

func Test_server_commentHandlerDelete(t *testing.T) {

	router := &Router{
		Commenter: &mockCommentStore{},
	}

	comment, err := router.Commenter.GetComment(uint(1))

	if err != nil {
		t.Errorf("Could not get mocked comment %v, got error %v", comment, err)
	}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name        string
		id          string
		statusCode  int
		messageBody *httpResponse
		errorBody   *httpResponse
	}{
		{
			name:       "Delete comment that exists in DB",
			id:         "1",
			statusCode: http.StatusOK,
			messageBody: &httpResponse{
				StatusCode:  http.StatusOK,
				Message:     http.StatusText(http.StatusOK),
				Description: "Comment successfully deleted.",
			},
			errorBody: nil,
		},
		{
			name:        "Delete comment that is not in db",
			id:          "1000",
			statusCode:  404,
			messageBody: nil,
			errorBody: &httpResponse{
				StatusCode:  404,
				Message:     http.StatusText(http.StatusNotFound),
				Description: "Could not find comment with id 1000.",
			},
		},
		{
			name:        "Get comment with invalid id parameter should cause Bad Request",
			id:          "-1",
			statusCode:  400,
			messageBody: nil,
			errorBody: &httpResponse{
				StatusCode:  400,
				Message:     http.StatusText(http.StatusBadRequest),
				Description: "Bad ID parameter -1.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", "/"+tt.id, nil)
			recorder := httptest.NewRecorder()
			muxRouter := router.Router()
			muxRouter.ServeHTTP(recorder, request)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected handler to respond with code %v, but got %v", tt.statusCode, recorder.Code)
			}

			if tt.messageBody != nil { // Expect regular body
				message := &httpResponse{}
				if err := json.NewDecoder(recorder.Body).Decode(message); err != nil {
					t.Errorf("Could not decode httpMessage body %v because of error %v", recorder.Body, err)
				}

				if !reflect.DeepEqual(message, tt.messageBody) {
					t.Errorf("Expected httpMessage %v, but was %v", tt.messageBody, message)
				}

			} else if tt.errorBody != nil { // Expect error body
				httpError := &httpResponse{}
				if err := json.NewDecoder(recorder.Body).Decode(httpError); err != nil {
					t.Errorf("Could not decode httpError body %v because of error %v", recorder.Body, err)
				}

				if !reflect.DeepEqual(httpError, tt.errorBody) {
					t.Errorf("Expected json error to be %v, but got %v", tt.errorBody, httpError)
				}
			}
		})
	}
}
