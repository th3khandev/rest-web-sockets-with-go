package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/th3khan/rest-web-sockets-with-go/models"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(url string) (*MySQLRepository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	return &MySQLRepository{
		db: db,
	}, nil
}

func (m *MySQLRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := m.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES (?, ?, ?)", user.ID, user.Email, user.Password)
	return err
}

func (m *MySQLRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = ?", id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user models.User
	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}
func (m *MySQLRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = ?", email)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user models.User
	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *MySQLRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := m.db.ExecContext(ctx, "INSERT INTO posts (id, title, content, user_id) VALUES (?, ?, ?, ?)", post.ID, post.Title, post.Content, post.UserID)
	return err
}

func (m *MySQLRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT id, title, content, user_id, created_at FROM posts WHERE id = ?", id)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var post models.Post
	for rows.Next() {
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err == nil {
			return &post, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &post, nil
}

func (m *MySQLRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := m.db.ExecContext(ctx, "UPDATE posts SET title = ?, content = ? WHERE id = ? AND user_id = ?", post.Title, post.Content, post.ID, post.UserID)
	return err
}

func (m *MySQLRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := m.db.ExecContext(ctx, "DELETE FROM posts WHERE id = ? AND user_id = ?", id, userId)
	return err
}

func (m *MySQLRepository) Close() error {
	return m.db.Close()
}
