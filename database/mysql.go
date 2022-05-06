package database

import (
	"context"
	"database/sql"
	"log"

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
	_, err := m.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
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

func (m *MySQLRepository) Close() error {
	return m.db.Close()
}
