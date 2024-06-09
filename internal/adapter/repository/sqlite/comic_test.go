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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
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
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetMaxId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"max(id)"}).AddRow(10)
	mock.ExpectQuery("select *").WillReturnRows(rows)
	id, err := cRepo.GetMaxId(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 10, id)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_Exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"url"}).AddRow("url")
	mock.ExpectQuery("select url from Comics").WillReturnRows(rows)
	exists, err := cRepo.Exists(context.Background(), 1)
	assert.True(t, exists)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_Exists_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	mock.ExpectQuery("select url from Comics").WillReturnError(sql.ErrNoRows)
	exists, err := cRepo.Exists(context.Background(), 1)
	assert.False(t, exists)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_Exists_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	mock.ExpectQuery("select url from Comics").WillReturnError(resErr)
	exists, err := cRepo.Exists(context.Background(), 1)
	assert.False(t, exists)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_Size(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(10)
	mock.ExpectQuery("select count *").WillReturnRows(rows)
	sz, err := cRepo.Size(context.Background())
	assert.Equal(t, 10, sz)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetUrls_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	mock.ExpectQuery("select id, url from Comics").WillReturnError(resErr)
	res, err := cRepo.GetUrls(context.Background())
	assert.Nil(t, res)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetUrls_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow("url", 2)
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	res, err := cRepo.GetUrls(context.Background())
	assert.Nil(t, res)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetUrls_RowError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow("url", 2).RowError(0, resErr)
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	res, err := cRepo.GetUrls(context.Background())
	assert.Nil(t, res)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetUrls(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow(1, "url").AddRow(2, "url2")
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	res, err := cRepo.GetUrls(context.Background())
	assert.Equal(t, map[int]string{1: "url", 2: "url2"}, res)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetIndex_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	mock.ExpectQuery("select id from Indexes where keyword").WithArgs("keyword").WillReturnError(resErr)
	res, err := cRepo.GetIndex(context.Background(), "keyword")
	assert.Nil(t, res)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetIndex_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow("id")
	mock.ExpectQuery("select id from Indexes where keyword").WithArgs("keyword").WillReturnRows(rows)
	res, err := cRepo.GetIndex(context.Background(), "keyword")
	assert.Nil(t, res)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetIndex_RowError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow("id").RowError(0, resErr)
	mock.ExpectQuery("select id from Indexes where keyword").WithArgs("keyword").WillReturnRows(rows)
	res, err := cRepo.GetIndex(context.Background(), "keyword")
	assert.Nil(t, res)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetIndex(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2)
	mock.ExpectQuery("select id from Indexes where keyword").WithArgs("keyword").WillReturnRows(rows)
	res, err := cRepo.GetIndex(context.Background(), "keyword")
	assert.Equal(t, []int{1, 2}, res)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetAll_QueryError1(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	mock.ExpectQuery("select id, url from Comics").WillReturnError(resErr)
	res, err := cRepo.GetAll(context.Background())
	assert.Nil(t, res)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetAll_ScanError1(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow("id", 2)
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	res, err := cRepo.GetAll(context.Background())
	assert.Nil(t, res)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetAll_QueryError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow(1, "url").AddRow(2, "url")
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	mock.ExpectQuery("select keyword from Indexes where id*").WillReturnError(resErr)
	res, err := cRepo.GetAll(context.Background())
	assert.Nil(t, res)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetAll_ScanError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow(1, "url").AddRow(2, "url")
	rows2 := sqlmock.NewRows([]string{"keyword"}).AddRow(5)
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	mock.ExpectQuery("select keyword from Indexes where id*").WillReturnRows(rows2)
	res, err := cRepo.GetAll(context.Background())
	assert.Nil(t, res)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetAll_RowError1(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow("id", 2).RowError(0, resErr)
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	res, err := cRepo.GetAll(context.Background())
	assert.Nil(t, res)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetAll_RowError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow(1, "url").AddRow(2, "url")
	rows2 := sqlmock.NewRows([]string{"keyword"}).AddRow(5).RowError(0, resErr)
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	mock.ExpectQuery("select keyword from Indexes where id*").WillReturnRows(rows2)
	res, err := cRepo.GetAll(context.Background())
	assert.Nil(t, res)
	assert.Equal(t, resErr, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestComicDB_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	cRepo := NewComicDB(db)

	rows := sqlmock.NewRows([]string{"id", "url"}).AddRow(1, "url")
	rows2 := sqlmock.NewRows([]string{"keyword"}).AddRow("appl").AddRow("day")
	mock.ExpectQuery("select id, url from Comics").WillReturnRows(rows)
	mock.ExpectQuery("select keyword from Indexes where id*").WithArgs(1).WillReturnRows(rows2)
	res, err := cRepo.GetAll(context.Background())
	assert.Equal(t, map[int]domain.Comic{1: domain.Comic{Url: "url", Keywords: []string{"appl", "day"}}}, res)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
