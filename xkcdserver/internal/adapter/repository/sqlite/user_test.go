package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var resErr = fmt.Errorf("error")

func TestUserDB_Add(t *testing.T) {
	user := domain.User{Name: "user", Password: "passwd", IsAdmin: true}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	uRepo := NewUserDB(db)
	mock.ExpectExec("INSERT INTO Users").WithArgs(user.Name, user.Password, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = uRepo.Add(context.Background(), user)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserDB_Add_Error(t *testing.T) {
	user := domain.User{Name: "user", Password: "passwd", IsAdmin: false}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	uRepo := NewUserDB(db)
	mock.ExpectExec("INSERT INTO Users").WithArgs(user.Name, user.Password, 0).WillReturnError(resErr)
	err = uRepo.Add(context.Background(), user)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserDB_GetByName(t *testing.T) {
	user := domain.User{Name: "user", Password: "passwd", IsAdmin: false}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	uRepo := NewUserDB(db)
	rows := sqlmock.NewRows([]string{"password", "is_admin"}).AddRow(user.Password, 0)
	mock.ExpectQuery("select password, is_admin from Users where name=").WithArgs(user.Name).WillReturnRows(rows)
	u, err := uRepo.GetByName(context.Background(), user.Name)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Equal(t, user, *u)
}

func TestUserDB_GetByName_Error(t *testing.T) {
	user := domain.User{Name: "user", Password: "passwd", IsAdmin: false}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	uRepo := NewUserDB(db)

	mock.ExpectQuery("select password, is_admin from Users where name=").WithArgs(user.Name).WillReturnError(resErr)
	u, err := uRepo.GetByName(context.Background(), user.Name)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Nil(t, u)
}

func TestUserDB_GetByName_NoRows(t *testing.T) {
	user := domain.User{Name: "user", Password: "passwd", IsAdmin: false}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	uRepo := NewUserDB(db)

	mock.ExpectQuery("select password, is_admin from Users where name=").WithArgs(user.Name).WillReturnError(sql.ErrNoRows)
	u, err := uRepo.GetByName(context.Background(), user.Name)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Nil(t, u)
}
