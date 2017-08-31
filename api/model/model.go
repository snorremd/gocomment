package model

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// CommentStore exposes common methods to create, get, update, and delete comments
type CommentStore interface {
	GetComments(string) ([]*Comment, error)
	GetComment(uint) (*Comment, error)
	CreateComment(*Comment) (*Comment, error)
	UpdateComment(*Comment) (*Comment, error)
	DeleteComment(*Comment) (*Comment, error)
}

// Comment represents a user comment
type Comment struct {
	ID        *uint      `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time `json:"createdAt" sql:"index"`
	UpdatedAt *time.Time `json:"updatedAt" sql:"index"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
	ParentID  uint       `json:"parentId"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Content   string     `json:"content"`
	Upvotes   int        `json:"upvotes"`
	Downvotes int        `json:"downvotes"`
	Status    string     `json:"status"`
	URL       string     `json:"url"`
}

// Migrate creates comment table using supplied db instance
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Comment{}).Error
}

// SqliteCommentStore implements a gorm based comment store
type SqliteCommentStore struct {
	DB *gorm.DB
}

// Validate checks if comment contains DB created fields
func Validate(comment *Comment) error {
	deniedFields := []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt"}

	v := reflect.ValueOf(*comment)

	props := make([]string, 0)

	for _, property := range deniedFields {
		if v.FieldByName(property).Pointer() != 0 {
			props = append(props, property)
		}
	}

	if len(props) > 0 {
		return errors.New("Following properties were not nil: " + strings.Join(props, ", "))
	}

	return nil
}

// GetComments fetches comments from database
func (c SqliteCommentStore) GetComments(url string) ([]*Comment, error) {
	comments := []*Comment{}
	return comments, c.DB.Where(&Comment{URL: url}).Find(&comments).Error
}

// GetComment fetches comment by id from database
func (c SqliteCommentStore) GetComment(id uint) (*Comment, error) {
	comment := Comment{}
	return &comment, c.DB.First(&comment, id).Error
}

// CreateComment inserts comment into database
func (c SqliteCommentStore) CreateComment(comment *Comment) (*Comment, error) {
	return comment, c.DB.Create(comment).Error
}

// UpdateComment updates selected comment
func (c SqliteCommentStore) UpdateComment(comment *Comment) (*Comment, error) {

	db := c.DB.Model(comment).Updates(comment)
	if db.Error != nil {
		return nil, db.Error
	} else if db.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	} else if err := db.Find(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment deletes selected comment
func (c SqliteCommentStore) DeleteComment(comment *Comment) (*Comment, error) {

	db := c.DB.Delete(comment)

	if db.Error != nil {
		return nil, db.Error
	} else if db.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return comment, nil
}
