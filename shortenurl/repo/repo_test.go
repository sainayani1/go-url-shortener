package repo

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"
	"url-shortener/shortenurl/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	return db, mock
}

func TestCreate(t *testing.T) {
	db, mock := newMock(t)
	query := "INSERT INTO short_urls (url, short_code, created_at, updated_at, access_count)VALUES(?, ?, ?, ?, ?)"
	currentTime := time.Now().UTC()
	testCases := []struct {
		description   string
		input         *model.ShortURL
		dbMock        interface{}
		output        int
		expectedError error
	}{
		{
			description: "failure case",
			input: &model.ShortURL{
				URL:       "https://www.example.com",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			dbMock: []interface{}{
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnError(sql.ErrConnDone),
			},
			output:        0,
			expectedError: sql.ErrConnDone,
		},
		{
			description: "success case",
			input: &model.ShortURL{
				URL:       "https://www.example.com",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			dbMock: []interface{}{
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnResult(sqlmock.NewResult(1, 0)),
			},
			output: 1,
		},
	}

	s := NewURLRepo(db)
	for i, tc := range testCases {
		output, err := s.Create(tc.input)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("Expected : %v\nGot : %v\n", tc.expectedError, err)
		}
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Test Case %v Failed. Expected %v Got %v", i+1, tc.output, output)
		}
	}
}

func TestUpdate(t *testing.T) {
	db, mock := newMock(t)
	query := "UPDATE short_urls SET url = ?, updated_at = ? WHERE short_code = ?"
	testCases := []struct {
		description   string
		code, newURL  string
		dbMock        interface{}
		expectedError error
	}{
		{
			description: "failure case",
			code:        "xsfg",
			newURL:      "https://www.example.co",
			dbMock: []interface{}{
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnError(sql.ErrConnDone),
			},
			expectedError: sql.ErrConnDone,
		},
		{
			description: "success case",
			code:        "xsfg",
			newURL:      "https://www.example.co",
			dbMock: []interface{}{
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnResult(sqlmock.NewResult(1, 0)),
			},
		},
	}

	s := NewURLRepo(db)
	for _, tc := range testCases {
		err := s.UpdateURL(tc.code, tc.newURL)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("Expected : %v\nGot : %v\n", tc.expectedError, err)
		}
	}
}

func TestDelete(t *testing.T) {
	db, mock := newMock(t)
	query := "DELETE FROM short_urls WHERE short_code = ?"
	testCases := []struct {
		description   string
		code          string
		dbMock        interface{}
		expectedError error
	}{
		{
			description: "failure case",
			code:        "xsfg",

			dbMock: []interface{}{
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnError(sql.ErrConnDone),
			},
			expectedError: sql.ErrConnDone,
		},
		{
			description: "success case",
			code:        "xsfg",
			dbMock: []interface{}{
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnResult(sqlmock.NewResult(1, 0)),
			},
		},
	}

	s := NewURLRepo(db)
	for _, tc := range testCases {
		err := s.DeleteByCode(tc.code)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("Expected : %v\nGot : %v\n", tc.expectedError, err)
		}
	}
}

func TestGetTrainersAndStudentCounts(t *testing.T) {
	db, mock := newMock(t)
	current := time.Now().UTC()

	query := "SELECT id, url, short_code, access_count, created_at, updated_at FROM short_urls WHERE short_code = ?"
	testCases := []struct {
		description   string
		code          string
		dbMock        interface{}
		output        *model.ShortURL
		expectedError error
	}{
		{
			description:        "failure case",
			code: "xyx",
			dbMock: []interface{}{
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(sqlmock.NewRows([]string{"id"})).WillReturnError(errors.New("sql error")),
			},
			output:        nil,
			expectedError: errors.New("sql error"),
		},

		{
			description: "success case",
			code:        "yz",
			output: &model.ShortURL{
				ID: 1,
				URL: "https://www.example.co",
				ShortCode: "yz",
				AccessCount: 1,
				CreatedAt: current,
				UpdatedAt: current,
			},
			dbMock: []interface{}{

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(sqlmock.NewRows([]string{"id", "url", "short_code", " access_count","created_at", "updated_at"}).AddRow(
					1, "https://www.example.co", "yz", 1, current, current,
				)),
			},
		},
	}

	s := NewURLRepo(db)
	for i, tc := range testCases {
		output, err := s.GetByCode(tc.code)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("Expected : %v\nGot : %v\n", tc.expectedError, err)
		}
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Test Case %v Failed. Expected %v Got %v", i+1, tc.output, output)
		}
	}
}
