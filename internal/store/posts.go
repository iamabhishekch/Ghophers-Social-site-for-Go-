package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"userid"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdat"`
	UpdatedAt string   `json:"updatedat"`
	// first_name and last_name can be added
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
    INSERT INTO posts (content, title, user_id, tags)
    VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
`
	err := s.db.QueryRowContext(
		ctx,
		query, 
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err!= nil{
		return err
	}
	
	return nil
}
