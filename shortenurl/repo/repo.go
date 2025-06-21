package repo

import (
	"database/sql"
	"time"
	"url-shortener/shortenurl/model"
)

type UrlRepo interface {
	Create(url *model.ShortURL) (int, error)
	GetByCode(code string) (*model.ShortURL, error)
	UpdateURL(code, newURL string) error
	DeleteByCode(code string) error
	IncrementAccessCount(code string) error
}

type UrlRepoStruct struct {
	db *sql.DB
}

func NewURLRepo(db *sql.DB) UrlRepoStruct {
	return UrlRepoStruct{db: db}
}

func (u *UrlRepoStruct) Create(url *model.ShortURL) (int, error) {
	resp, err := u.db.Exec("INSERT INTO short_urls (url, short_code, created_at, updated_at, access_count)VALUES(?, ?, ?, ?, ?)", url.URL, url.ShortCode, url.CreatedAt, url.UpdatedAt, url.AccessCount)
	if err != nil {
		return 0, err
	}

	lastInsertedVal, err := resp.LastInsertId()

	return int(lastInsertedVal), err
}

func (u *UrlRepoStruct) GetByCode(code string) (*model.ShortURL, error) {
	var resp model.ShortURL
	err := u.db.QueryRow("SELECT id, url, short_code, access_count, created_at, updated_at FROM short_urls WHERE short_code = ?", code).
		Scan(&resp.ID, &resp.URL, &resp.ShortCode, &resp.AccessCount, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (u *UrlRepoStruct) UpdateURL(code, newURL string) error {
	_, err := u.db.Exec("UPDATE short_urls SET url = ?, updated_at = ? WHERE short_code = ?", newURL, time.Now(), code)

	return err
}

func (u *UrlRepoStruct) DeleteByCode(code string) error {
	_, err := u.db.Exec("DELETE FROM short_urls WHERE short_code = ?", code)
	return err
}

func (u *UrlRepoStruct) IncrementAccessCount(code string) error {
	_, err := u.db.Exec("UPDATE short_urls SET access_count = access_count + 1, updated_at = ? WHERE short_code = ?", time.Now(), code)
	return err
}
