package models

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Excerpt   string    `json:"excerpt"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostStore struct {
	DB *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{DB: db}
}

// GenerateSlug creates a URL-safe slug from a title
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// GetPublishedPosts returns all published posts ordered by creation date
func (s *PostStore) GetPublishedPosts() ([]Post, error) {
	rows, err := s.DB.Query(`
		SELECT id, title, slug, excerpt, content, published, created_at, updated_at
		FROM posts
		WHERE published = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query published posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// GetAllPosts returns all posts (for admin) ordered by creation date
func (s *PostStore) GetAllPosts() ([]Post, error) {
	rows, err := s.DB.Query(`
		SELECT id, title, slug, excerpt, content, published, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query all posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// GetPostBySlug returns a single post by its slug
func (s *PostStore) GetPostBySlug(slug string) (*Post, error) {
	var p Post
	err := s.DB.QueryRow(`
		SELECT id, title, slug, excerpt, content, published, created_at, updated_at
		FROM posts
		WHERE slug = $1
	`, slug).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query post by slug: %w", err)
	}
	return &p, nil
}

// GetPostByID returns a single post by its ID
func (s *PostStore) GetPostByID(id int) (*Post, error) {
	var p Post
	err := s.DB.QueryRow(`
		SELECT id, title, slug, excerpt, content, published, created_at, updated_at
		FROM posts
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query post by id: %w", err)
	}
	return &p, nil
}

// CreatePost inserts a new post
func (s *PostStore) CreatePost(p *Post) error {
	if p.Slug == "" {
		p.Slug = GenerateSlug(p.Title)
	}
	err := s.DB.QueryRow(`
		INSERT INTO posts (title, slug, excerpt, content, published)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, p.Title, p.Slug, p.Excerpt, p.Content, p.Published).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert post: %w", err)
	}
	return nil
}

// UpdatePost updates an existing post
func (s *PostStore) UpdatePost(p *Post) error {
	if p.Slug == "" {
		p.Slug = GenerateSlug(p.Title)
	}
	_, err := s.DB.Exec(`
		UPDATE posts
		SET title = $1, slug = $2, excerpt = $3, content = $4, published = $5, updated_at = NOW()
		WHERE id = $6
	`, p.Title, p.Slug, p.Excerpt, p.Content, p.Published, p.ID)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	return nil
}

// DeletePost deletes a post by ID
func (s *PostStore) DeletePost(id int) error {
	_, err := s.DB.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}
