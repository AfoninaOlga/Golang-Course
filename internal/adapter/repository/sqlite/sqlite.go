package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
)

type ComicDB struct {
	db *sql.DB
}

func New(db *sql.DB) *ComicDB {
	return &ComicDB{db: db}
}

func (cdb *ComicDB) GetAll(ctx context.Context) (map[int]domain.Comic, error) {
	comicMap := make(map[int]domain.Comic)
	comics, err := cdb.db.QueryContext(ctx, "select id, url from Comics")
	if err != nil {
		return nil, err
	}
	defer comics.Close()
	for comics.Next() {
		var (
			id       int
			url      string
			keywords []string
			word     string
		)
		if err = comics.Scan(&id, &url); err != nil {
			return nil, err
		}
		words, err := cdb.db.QueryContext(ctx, "select keyword from Indexes where id=? ", id)
		if err != nil {
			return nil, err
		}
		defer words.Close()
		for words.Next() {
			err = words.Scan(&word)
			if err != nil {
				return nil, err
			}
			keywords = append(keywords, word)
		}
		comicMap[id] = domain.Comic{Url: url, Keywords: keywords}
	}
	return comicMap, nil
}

func (cdb *ComicDB) GetIndex(ctx context.Context, keyword string) ([]int, error) {
	var res []int
	ids, err := cdb.db.QueryContext(ctx, "select id from Indexes where keyword=?", keyword)
	if err != nil {
		return nil, err
	}
	defer ids.Close()
	for ids.Next() {
		var id int
		if err = ids.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (cdb *ComicDB) AddComic(ctx context.Context, id int, c domain.Comic) error {
	tx, err := cdb.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	insertComic, err := tx.Prepare("INSERT INTO Comics(id, url) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer insertComic.Close()

	insertKeyword, err := tx.Prepare("INSERT INTO Keywords(word) VALUES (?)")
	if err != nil {
		return err
	}
	defer insertKeyword.Close()

	insertComicKeyword, err := tx.Prepare("INSERT INTO ComicsKeywords(comic_id, keyword_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer insertComicKeyword.Close()

	getKeyworId, err := tx.Prepare("select id from Keywords where word=?")
	if err != nil {
		return err
	}
	defer insertComicKeyword.Close()

	_, err = insertComic.Exec(id, c.Url)
	if err != nil {
		return err
	}

	for _, word := range c.Keywords {
		var kId int64
		err = getKeyworId.QueryRow(word).Scan(&kId)
		if errors.Is(err, sql.ErrNoRows) {
			res, err := insertKeyword.Exec(word)
			if err != nil {
				return err
			}
			kId, err = res.LastInsertId()
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		_, err = insertComicKeyword.Exec(id, kId)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func (cdb *ComicDB) GetMaxId(ctx context.Context) (int, error) {
	var id int
	err := cdb.db.QueryRowContext(ctx, "select max(id) from Comics").Scan(&id)
	return id, err
}

func (cdb *ComicDB) Exists(ctx context.Context, id int) (bool, error) {
	var url string
	err := cdb.db.QueryRowContext(ctx, "select url from Comics where id=?", id).Scan(&url)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func (cdb *ComicDB) Size(ctx context.Context) (int, error) {
	var sz int
	err := cdb.db.QueryRowContext(ctx, "select count(*) from Comics").Scan(&sz)
	return sz, err
}

func (cdb *ComicDB) GetUrls(ctx context.Context) (map[int]string, error) {
	res := map[int]string{}
	comics, err := cdb.db.QueryContext(ctx, "select id, url from Comics")
	if err != nil {
		return nil, err
	}
	defer comics.Close()
	for comics.Next() {
		var (
			id  int
			url string
		)
		if err = comics.Scan(&id, &url); err != nil {
			return nil, err
		}
		res[id] = url
	}
	return res, nil
}
