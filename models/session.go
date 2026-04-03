package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
)

type SessionStore struct {
	DB *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{DB: db}
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *SessionStore) CreateSession(userID int) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	// Session valid for 24 hours
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = s.DB.Exec(`
		INSERT INTO sessions (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("insert session: %w", err)
	}

	return token, nil
}

func (s *SessionStore) ValidateSession(token string) (*User, error) {
	var u User
	var expiresAt time.Time

	err := s.DB.QueryRow(`
		SELECT u.id, u.username, u.password_hash, u.created_at, s.expires_at
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.token = $1
	`, token).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt, &expiresAt)

	if err != nil {
		return nil, err // Could be sql.ErrNoRows
	}

	if time.Now().After(expiresAt) {
		s.DeleteSession(token)
		return nil, fmt.Errorf("session expired")
	}

	return &u, nil
}

func (s *SessionStore) DeleteSession(token string) error {
	_, err := s.DB.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}
