package service

import (
	"errors"
	"reflect"
	"testing"
	"time"
	"url-shortener/shortenurl/mocks"
	"url-shortener/shortenurl/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {
	mockService := new(mocks.UrlRepo)
	now := time.Now().UTC()

	type testCase struct {
		desc           string
		input          string
		mockSetup      func()
		expectedOutput *model.ShortURL
		expectedError  error
	}

	testCases := []testCase{
		{
			desc:  "success case",
			input: "https://www.example.co",
			mockSetup: func() {
				mockService.On("Create", mock.Anything).Return(1, nil)
			},
			expectedOutput: &model.ShortURL{ID: 1, URL: "https://www.example.co", ShortCode: "bYFXNv", CreatedAt: now, UpdatedAt: now},
		},

		{
			desc:  "failure case",
			input: "https://www.example.co",
			mockSetup: func() {
				mockService = new(mocks.UrlRepo)
				mockService.On("Create", mock.Anything).Return(0, errors.New("Please try again"))
			},
			expectedError: errors.New("Please try again"),
		},
	}

	for i, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()

			str := NewURLService(mockService)

			output, err := str.Create(tc.input)

			assert.Equal(t, tc.expectedError, err)

			if tc.expectedOutput != nil || output != nil {
				if !reflect.DeepEqual(tc.expectedOutput.ID, output.ID) {
					t.Errorf("Test Case %v Failed. Expected %v Got %v", i+1, tc.expectedOutput, output)
				}
			}

		})
	}
}

func TestServiceGet(t *testing.T) {
	mockService := new(mocks.UrlRepo)
	now := time.Now().UTC()

	type testCase struct {
		desc           string
		input,url          string
		mockSetup      func()
		expectedOutput *model.ShortURL
		expectedError  error
	}

	testCases := []testCase{
		{
			desc:  "success case",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService.On("GetByCode", mock.Anything).Return(&model.ShortURL{ID: 1, URL: "https://www.example.co", ShortCode: "bYFXNv", CreatedAt: now, UpdatedAt: now},nil)
				mockService.On("IncrementAccessCount", mock.Anything).Return(nil)
			},
			expectedOutput: &model.ShortURL{ID: 1, URL: "https://www.example.co", ShortCode: "bYFXNv", CreatedAt: now, UpdatedAt: now,AccessCount: 1},
		},

		{
			desc:  "failure case: error from store layer details",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService = new(mocks.UrlRepo)
				mockService.On("GetByCode", mock.Anything).Return(nil, errors.New("Please try again"))
			},
			expectedError: errors.New(errShortURLNotFound.Error()+"Please try again"),
		},

		{
			desc:  "failure case: error from store layer details",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService = new(mocks.UrlRepo)
				mockService.On("GetByCode", mock.Anything).Return(&model.ShortURL{ID: 1, URL: "https://www.example.co", ShortCode: "bYFXNv", CreatedAt: now, UpdatedAt: now},nil)
				mockService.On("IncrementAccessCount", mock.Anything).Return(errors.New("Please try again"))
			},
			expectedError: errors.New(errAccessUpdate.Error() + "Please try again"),
		},
	}

	for i, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()

			str := NewURLService(mockService)

			output, err := str.Get(tc.input)

			assert.Equal(t, tc.expectedError, err)

		
				if tc.expectedOutput != nil || output != nil {
				if !reflect.DeepEqual(tc.expectedOutput.ID, output.ID) {
					t.Errorf("Test Case %v Failed. Expected %v Got %v", i+1, tc.expectedOutput, output)
				}
			}

		})
	}
}

func TestServiceUpdate(t *testing.T) {
	mockService := new(mocks.UrlRepo)
	now := time.Now().UTC()

	type testCase struct {
		desc           string
		input,url          string
		mockSetup      func()
		expectedOutput *model.ShortURL
		expectedError  error
	}

	testCases := []testCase{
		{
			desc:  "success case",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService.On("UpdateURL", mock.Anything, mock.Anything).Return(nil)
				mockService.On("GetByCode", mock.Anything).Return(&model.ShortURL{ID: 1, URL: "https://www.example.co", ShortCode: "bYFXNv", CreatedAt: now, UpdatedAt: now},nil)
				mockService.On("IncrementAccessCount", mock.Anything).Return(nil)
			},
			expectedOutput: &model.ShortURL{ID: 1, URL: "https://www.example.co", ShortCode: "bYFXNv", CreatedAt: now, UpdatedAt: now},
		},

		{
			desc:  "failure case: error from store layer details",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService = new(mocks.UrlRepo)
				mockService.On("UpdateURL", mock.Anything, mock.Anything).Return(errors.New("Please try again"))
			},
			expectedError: errors.New("Please try again"),
		},

		{
			desc:  "failure case: error from store layer details",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService.On("UpdateURL", mock.Anything, mock.Anything).Return(nil)
				mockService.On("GetByCode", mock.Anything).Return(nil, errors.New("Please try again"))
			},
			expectedError: errors.New("Please try again"),
		},

		{
			desc:  "failure case: error from store layer details",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService.On("UpdateURL", mock.Anything, mock.Anything).Return(nil)
				mockService.On("GetByCode", mock.Anything).Return(&model.ShortURL{ID: 1, URL: "https://www.example.co", ShortCode: "bYFXNv", CreatedAt: now, UpdatedAt: now},nil)
				mockService.On("IncrementAccessCount", mock.Anything).Return(errors.New("Please try again"))
			},
			expectedError: errors.New("Please try again"),
		},
	}

	for i, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()

			str := NewURLService(mockService)

			output, err := str.Update(tc.input,tc.url)

			assert.Equal(t, tc.expectedError, err)

			if tc.expectedOutput != nil || output != nil {
				if !reflect.DeepEqual(tc.expectedOutput.ID, output.ID) {
					t.Errorf("Test Case %v Failed. Expected %v Got %v", i+1, tc.expectedOutput, output)
				}
			}

		})
	}
}

func TestServiceDelete(t *testing.T) {
	mockService := new(mocks.UrlRepo)

	type testCase struct {
		desc           string
		input,url          string
		mockSetup      func()
		expectedError  error
	}

	testCases := []testCase{
		{
			desc:  "success case",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService.On("DeleteByCode", mock.Anything).Return(nil)
			},
		},

		{
			desc:  "failure case: error from store layer details",
			input: "cccc",
			url: "https://www.example.co",
			mockSetup: func() {
				mockService = new(mocks.UrlRepo)
				mockService.On("DeleteByCode",  mock.Anything).Return(errors.New("Please try again"))
			},
			expectedError: errors.New("Please try again"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()

			str := NewURLService(mockService)

			err := str.Delete(tc.input)

			assert.Equal(t, tc.expectedError, err)

			

		})
	}
}
