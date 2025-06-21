package service

import (
	"errors"
	"time"
	"url-shortener/common"
	"url-shortener/shortenurl/model"
	"url-shortener/shortenurl/repo"
)

var (
	errShortURLNotFound = errors.New("short URL not found : ")
	errAccessUpdate     = errors.New("access count update failed : ")
)

type IURLService interface {
	Create(url string) (*model.ShortURL, error)
	Get(code string) (*model.ShortURL, error)
	Update(code, newURL string) (*model.ShortURL, error)
	Delete(code string) error
}

type urlService struct {
	repo repo.UrlRepo
}

func NewURLService(repo repo.UrlRepo) IURLService {
	return &urlService{repo: repo}
}

func (svc *urlService) Create(url string) (*model.ShortURL, error) {
	now := time.Now().UTC()
	request := model.ShortURL{
		URL:         url,
		ShortCode:   common.GenerateShortCode(6),
		AccessCount: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	lastInsID, err := svc.repo.Create(&request)
	if err != nil {
		return nil, err
	}

	request.ID = lastInsID
	return &request, nil
}

func (s *urlService) Get(code string) (*model.ShortURL, error) {
	url, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, errors.New(errShortURLNotFound.Error() + err.Error())
	}

	if err = s.repo.IncrementAccessCount(code); err != nil {
		return nil, errors.New(errAccessUpdate.Error() + err.Error())
	}

	url.AccessCount++

	return url, nil
}

func (svc *urlService) Update(code, newURL string) (*model.ShortURL, error) {
	err := svc.repo.UpdateURL(code, newURL)
	if err != nil {
		return nil, err
	}

	return svc.Get(code)
}

func (svc *urlService) Delete(code string) error {
	return svc.repo.DeleteByCode(code)
}
