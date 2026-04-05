package models

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Excerpt    string    `json:"excerpt"`
	Content    string    `json:"content"`
	ImageURL   string    `json:"image_url"`
	Published  bool      `json:"published"`
	AuthorID   int       `json:"author_id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
		SELECT p.id, p.title, p.slug, p.excerpt, p.content, p.image_url, p.published, p.created_at, p.updated_at, COALESCE(u.username, 'Unknown'), COALESCE(p.author_id, 0)
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.published = true
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query published posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.ImageURL, &p.Published, &p.CreatedAt, &p.UpdatedAt, &p.AuthorName, &p.AuthorID); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// GetAllPosts returns all posts (for admin) ordered by creation date
func (s *PostStore) GetAllPosts() ([]Post, error) {
	rows, err := s.DB.Query(`
		SELECT p.id, p.title, p.slug, p.excerpt, p.content, p.image_url, p.published, p.created_at, p.updated_at, COALESCE(u.username, 'Unknown'), COALESCE(p.author_id, 0)
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query all posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.ImageURL, &p.Published, &p.CreatedAt, &p.UpdatedAt, &p.AuthorName, &p.AuthorID); err != nil {
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
		SELECT p.id, p.title, p.slug, p.excerpt, p.content, p.image_url, p.published, p.created_at, p.updated_at, COALESCE(u.username, 'Unknown'), COALESCE(p.author_id, 0)
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.slug = $1
	`, slug).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.ImageURL, &p.Published, &p.CreatedAt, &p.UpdatedAt, &p.AuthorName, &p.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("query post by slug: %w", err)
	}
	return &p, nil
}

// GetPostByID returns a single post by its ID
func (s *PostStore) GetPostByID(id int) (*Post, error) {
	var p Post
	err := s.DB.QueryRow(`
		SELECT p.id, p.title, p.slug, p.excerpt, p.content, p.image_url, p.published, p.created_at, p.updated_at, COALESCE(u.username, 'Unknown'), COALESCE(p.author_id, 0)
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.ImageURL, &p.Published, &p.CreatedAt, &p.UpdatedAt, &p.AuthorName, &p.AuthorID)
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
	var authorID interface{} = p.AuthorID
	if p.AuthorID == 0 {
		authorID = nil
	}
	err := s.DB.QueryRow(`
		INSERT INTO posts (title, slug, excerpt, content, image_url, published, author_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`, p.Title, p.Slug, p.Excerpt, p.Content, p.ImageURL, p.Published, authorID).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
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
		SET title = $1, slug = $2, excerpt = $3, content = $4, image_url = $5, published = $6, updated_at = NOW()
		WHERE id = $7
	`, p.Title, p.Slug, p.Excerpt, p.Content, p.ImageURL, p.Published, p.ID)
	// Note: author_id is intentionally not updated here per design choice.
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
