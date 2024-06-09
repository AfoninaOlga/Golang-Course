package sqlite

import (
	"context"
	"database/sql"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComicDB_AddComic(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*")
	mock.ExpectPrepare("INSERT INTO ComicsKeywords(comic_id, keyword_id)*")
	mock.ExpectPrepare("select id from Keywords where word*")
	mock.ExpectExec("INSERT INTO Comics(id, url)*").WithArgs(1, comic.Url).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("select id from Keywords where word*").WithArgs(comic.Keywords[0]).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO Keywords(word)*").WithArgs(comic.Keywords[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO ComicsKeywords(comic_id, keyword_id)*").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorBegin(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin().WillReturnError(resErr)
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorPrepare1(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*").WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorPerpare2(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*").WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorPerpare3(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*")
	mock.ExpectPrepare("INSERT INTO ComicsKeywords(comic_id, keyword_id)*").WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorPerpare4(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*")
	mock.ExpectPrepare("INSERT INTO ComicsKeywords(comic_id, keyword_id)*")
	mock.ExpectPrepare("select id from Keywords where word*").WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorExec1(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*")
	mock.ExpectPrepare("INSERT INTO ComicsKeywords(comic_id, keyword_id)*")
	mock.ExpectPrepare("select id from Keywords where word*")
	mock.ExpectExec("INSERT INTO Comics(id, url)*").WithArgs(1, comic.Url).WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorExec2(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*")
	mock.ExpectPrepare("INSERT INTO ComicsKeywords(comic_id, keyword_id)*")
	mock.ExpectPrepare("select id from Keywords where word*")
	mock.ExpectExec("INSERT INTO Comics(id, url)*").WithArgs(1, comic.Url).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("select id from Keywords where word*").WithArgs(comic.Keywords[0]).WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorExec3(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*")
	mock.ExpectPrepare("INSERT INTO ComicsKeywords(comic_id, keyword_id)*")
	mock.ExpectPrepare("select id from Keywords where word*")
	mock.ExpectExec("INSERT INTO Comics(id, url)*").WithArgs(1, comic.Url).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("select id from Keywords where word*").WithArgs(comic.Keywords[0]).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO Keywords(word)*").WithArgs(comic.Keywords[0]).WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_AddComic_ErrorExec4(t *testing.T) {
	comic := domain.Comic{Url: "url", Keywords: []string{"apple"}}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO Comics(id, url)*")
	mock.ExpectPrepare("INSERT INTO Keywords(word)*")
	mock.ExpectPrepare("INSERT INTO ComicsKeywords(comic_id, keyword_id)*")
	mock.ExpectPrepare("select id from Keywords where word*")
	mock.ExpectExec("INSERT INTO Comics(id, url)*").WithArgs(1, comic.Url).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("select id from Keywords where word*").WithArgs(comic.Keywords[0]).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO Keywords(word)*").WithArgs(comic.Keywords[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO ComicsKeywords(comic_id, keyword_id)*").WithArgs(1, 1).WillReturnError(resErr)
	mock.ExpectRollback()
	err = cRepo.AddComic(context.Background(), 1, comic)
	assert.Error(t, err, resErr)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

//func TestComicDB_GetMaxId(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	assert.NoError(t, err)
//	defer db.Close()
//	cRepo := NewComicDB(db)
//
//	rows := sqlmock.NewRows([]string{"max(id)"}).AddRow(10)
//	mock.ExpectQuery("select max(id) from Comics").WillReturnRows(rows)
//	id, err := cRepo.GetMaxId(context.Background())
//	assert.NoError(t, err)
//	assert.Equal(t, 10, id)
//}
