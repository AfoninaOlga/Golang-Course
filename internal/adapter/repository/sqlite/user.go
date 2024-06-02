package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	_ "github.com/mattn/go-sqlite3"
)

type UserDB struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{db: db}
}

func (udb *UserDB) Add(ctx context.Context, u domain.User) error {
	var admin int
	if u.IsAdmin {
		admin = 1
	}
	tx, err := udb.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	insertComic, err := tx.Prepare("INSERT INTO Users(name, password, is_admin) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer insertComic.Close()
	_, err = insertComic.Exec(u.Name, u.Password, admin)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (udb *UserDB) GetByName(ctx context.Context, name string) (*domain.User, error) {
	var (
		passwd  string
		isAdmin int
	)
	err := udb.db.QueryRowContext(ctx, "select password, is_admin from Users where name=?", name).Scan(&passwd, &isAdmin)

	if err == nil {
		return &domain.User{Name: name, Password: passwd, IsAdmin: isAdmin == 1}, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}
