package model

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/snorremd/gocomment/api/db"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func dbname() string {
	return uuid.New().String() + ".db"
}

func setupDB(t *testing.T, db *gorm.DB) {

	if err := db.DropTableIfExists("comments").Error; err != nil {
		t.FailNow()
	}

	if err := Migrate(db); err != nil {
		t.FailNow()
	}
}

func TestMigrate(t *testing.T) {

	dbname := dbname()
	db, err := db.DB(dbname)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}
	defer db.Close()
	defer os.Remove(dbname)

	tests := []struct {
		name string
	}{
		{
			name: "Successfully create table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Migrate(db)
			if err != nil {
				t.Errorf("Database migration failed.")
			}
		})
	}
}

func TestGetComments(t *testing.T) {
	dbname := dbname()
	db, err := db.DB(dbname)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}
	defer db.Close()
	defer os.Remove(dbname)
	setupDB(t, db)

	commenter := &SqliteCommentStore{DB: db}

	comment1, _ := commenter.CreateComment(&Comment{
		Content:   "Some content all right",
		Upvotes:   0,
		Downvotes: 0,
		Status:    "Approved",
		URL:       "http://example.com/post/1",
	})

	comment2, _ := commenter.CreateComment(&Comment{
		Content:   "More content all right",
		Upvotes:   0,
		Downvotes: 0,
		Status:    "Approved",
		URL:       "http://example.com/post/2",
	})

	comment3, _ := commenter.CreateComment(&Comment{
		Content:   "More content for post 2 all right",
		Upvotes:   0,
		Downvotes: 0,
		Status:    "Approved",
		URL:       "http://example.com/post/2",
	})

	commentsPost1 := []*Comment{comment1}
	commentsPost2 := []*Comment{comment2, comment3}

	tests := []struct {
		name     string
		url      string
		comments []*Comment
		wantErr  bool
	}{
		{
			name:     "Fetch comments for http://example.com/post/1",
			url:      "http://example.com/post/1",
			comments: commentsPost1,
			wantErr:  false,
		},
		{
			name:     "Fetch comments for http://example.com/post/2",
			url:      "http://example.com/post/2",
			comments: commentsPost2,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if comments, err := commenter.GetComments(tt.url); (err != nil) != tt.wantErr {
				t.Errorf("GetComments() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(comments) != len(tt.comments) {
				t.Errorf("GetComments() Expected to find %v comments, found %v", len(tt.comments), len(comments))
			}
		})
	}
}

func Test_createComment(t *testing.T) {

	dbname := dbname()
	db, err := db.DB(dbname)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}
	defer db.Close()
	defer os.Remove(dbname)
	setupDB(t, db)

	commenter := &SqliteCommentStore{DB: db}

	tests := []struct {
		name    string
		comment *Comment
		wantErr bool
	}{
		{
			name: "Create comment in empty db",
			comment: &Comment{
				Content:   "Some content all right",
				Upvotes:   0,
				Downvotes: 0,
				Status:    "Approved",
			},
			wantErr: false,
		},
		{
			name: "Create comment number two",
			comment: &Comment{
				Content:   "Some more content all right",
				Upvotes:   0,
				Downvotes: 0,
				Status:    "Disapproved",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := commenter.CreateComment(tt.comment); (err != nil) != tt.wantErr {
				t.Errorf("createComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getComment(t *testing.T) {

	dbname := dbname()
	db, err := db.DB(dbname)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}
	defer db.Close()
	defer os.Remove(dbname)
	setupDB(t, db)

	commenter := &SqliteCommentStore{DB: db}

	comment1, _ := commenter.CreateComment(&Comment{
		Content:   "Some content all right",
		Upvotes:   0,
		Downvotes: 0,
		Status:    "Approved",
	})

	comment2, _ := commenter.CreateComment(&Comment{
		Content:   "More content all right",
		Upvotes:   0,
		Downvotes: 0,
		Status:    "Approved",
	})

	tests := []struct {
		name    string
		comment *Comment
		wantErr bool
	}{
		{
			name:    "Successfully get comment1",
			comment: comment1,
			wantErr: false,
		},
		{
			name:    "Successfully get comment2",
			comment: comment2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if comment, err := commenter.GetComment(*tt.comment.ID); (err != nil) != tt.wantErr {
				t.Errorf("getComments() error = %v, wantErr %v", err, tt.wantErr)
			} else if *comment.ID != *tt.comment.ID {
				t.Errorf("getComments() wanted id = %v, but got id = %v", tt.comment.ID, comment.ID)
			}
		})
	}
}

func Test_updateComment(t *testing.T) {

	dbname := dbname()
	db, err := db.DB(dbname)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}
	defer db.Close()
	defer os.Remove(dbname)
	setupDB(t, db)

	commenter := &SqliteCommentStore{DB: db}

	comment1, _ := commenter.CreateComment(&Comment{
		Content:   "Some content all right",
		Upvotes:   0,
		Downvotes: 0,
		Status:    "Approved",
	})

	comment2Id := uint(1000)
	comment2 := &Comment{
		ID:      &comment2Id,
		Content: "Should not exist",
	}

	tests := []struct {
		name    string
		comment *Comment
		wantErr bool
	}{
		{
			name:    "Successfully update comment1",
			comment: comment1,
			wantErr: false,
		},
		{
			name:    "Update comment that does not exist",
			comment: comment2,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if comment, err := commenter.UpdateComment(tt.comment); (err != nil) != tt.wantErr {
				t.Errorf("updateComment() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.wantErr == false && comment.UpdatedAt == comment.CreatedAt {
				t.Errorf("UpdatedAt %v equals CreatedAt %v", comment.UpdatedAt, comment.CreatedAt)
			} else if tt.wantErr == false && !reflect.DeepEqual(comment, comment) {
				t.Errorf("Updated comment %v not equal to comment %v.", comment, tt.comment)
			}
		})
	}
}

func Test_deleteComment(t *testing.T) {

	dbname := dbname()
	db, err := db.DB(dbname)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}
	defer db.Close()
	defer os.Remove(dbname)
	setupDB(t, db)

	commenter := &SqliteCommentStore{DB: db}

	comment1, _ := commenter.CreateComment(&Comment{
		Content:   "Some content all right",
		Upvotes:   0,
		Downvotes: 0,
		Status:    "Approved",
	})

	tests := []struct {
		name    string
		comment *Comment
		wantErr bool
	}{
		{
			name:    "Successfully delete comment1",
			comment: comment1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := commenter.DeleteComment(tt.comment); (err != nil) != tt.wantErr {
				t.Errorf("deleteComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
