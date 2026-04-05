package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Comment struct {
	ID          int       `json:"id"`
	PostID      int       `json:"post_id"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Content     string    `json:"content"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Post relation for admin display
	PostTitle   string    `json:"post_title"` 
}

type CommentStore struct {
	DB *sql.DB
}

func NewCommentStore(db *sql.DB) *CommentStore {
	return &CommentStore{DB: db}
}

// GetCommentsByPostID gets all published comments for a specific post
func (s *CommentStore) GetCommentsByPostID(postID int) ([]Comment, error) {
	rows, err := s.DB.Query(`
		SELECT id, post_id, author_name, author_email, content, is_published, created_at
		FROM comments
		WHERE post_id = $1 AND is_published = true
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("query published comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorName, &c.AuthorEmail, &c.Content, &c.IsPublished, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// GetAllComments gets all comments for the admin moderation page
func (s *CommentStore) GetAllComments() ([]Comment, error) {
	rows, err := s.DB.Query(`
		SELECT c.id, c.post_id, c.author_name, c.author_email, c.content, c.is_published, c.created_at, p.title
		FROM comments c
		JOIN posts p ON c.post_id = p.id
		ORDER BY c.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query all comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorName, &c.AuthorEmail, &c.Content, &c.IsPublished, &c.CreatedAt, &c.PostTitle); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// CreateComment inserts a new comment
func (s *CommentStore) CreateComment(c *Comment) error {
	err := s.DB.QueryRow(`
		INSERT INTO comments (post_id, author_name, author_email, content, is_published)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, c.PostID, c.AuthorName, c.AuthorEmail, c.Content, c.IsPublished).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert comment: %w", err)
	}
	return nil
}

// PublishComment sets is_published to true
func (s *CommentStore) PublishComment(id int) error {
	_, err := s.DB.Exec(`UPDATE comments SET is_published = true WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("publish comment: %w", err)
	}
	return nil
}

// DeleteComment removes a comment
func (s *CommentStore) DeleteComment(id int) error {
	_, err := s.DB.Exec(`DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	return nil
}
