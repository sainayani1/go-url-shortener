package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/shortenurl/mocks"
	"url-shortener/shortenurl/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerCreate(t *testing.T) {
	mockService := new(mocks.IURLService)

	type testCase struct {
		desc           string
		body           string
		mockSetup      func()
		expectedStatus int
	}

	testCases := []testCase{
		{
			desc: "success case",
			body: `{"url": "https://www.example.co"}`,
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Create", mock.Anything).Return(&model.ShortURL{ID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			desc: "failure case",
			body: `{"url": 1}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},

		{
			desc: "failure case",
			body: `{"url": ""}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},

		{
			desc: "failure case",
			body: `{"url": "https://www.example.co"}`,
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Create", mock.Anything).Return(nil, errors.New("Please try again"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()
			h := URLHandler{service: mockService}
			handler := NewHandler(h.service)
			req, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewReader([]byte(tc.body)))
			assert.NoError(t, err)
			rec := httptest.NewRecorder()
			handler.Create(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rec.Code)
			}

		})
	}
}

func muxSetVars(req *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(req, vars)
}

func TestHandlerUpdate(t *testing.T) {
	mockService := new(mocks.IURLService)

	type testCase struct {
		desc           string
		body, code     string
		mockSetup      func()
		expectedStatus int
	}

	testCases := []testCase{
		{
			desc: "success case",
			body: `{"url": "https://www.example.co"}`,
			code: "xyx",
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Update", mock.Anything, mock.Anything).Return(&model.ShortURL{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			desc: "failure case",
			body: `{"url": 1}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},

		{
			desc: "failure case",
			body: `{"url": ""}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			desc: "failure case",
			body: `{"url": "https://www.example.co"}`,
			code: "xyx",
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("Please try again"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()
			h := URLHandler{service: mockService}
			handler := NewHandler(h.service)
			req, err := http.NewRequest(http.MethodPut, "/shorten/"+tc.code, bytes.NewReader([]byte(tc.body)))
			req = muxSetVars(req, map[string]string{"code": tc.code})
			assert.NoError(t, err)
			rec := httptest.NewRecorder()
			handler.Update(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rec.Code)
			}

		})
	}
}

func TestHandlerDelete(t *testing.T) {
	mockService := new(mocks.IURLService)

	type testCase struct {
		desc           string
		body, code     string
		mockSetup      func()
		expectedStatus int
	}

	testCases := []testCase{
		{
			desc: "success case",
			body: `{"url": "https://www.example.co"}`,
			code: "xyx",
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Delete", mock.Anything).Return( nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			desc: "failure case",
			body: `{"url": "https://www.example.co"}`,
			code: "xyx",
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Delete", mock.Anything, mock.Anything).Return(errors.New("Please try again"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()
			h := URLHandler{service: mockService}
			handler := NewHandler(h.service)
			req, err := http.NewRequest(http.MethodDelete, "/shorten/"+tc.code, bytes.NewReader([]byte(tc.body)))
			req = muxSetVars(req, map[string]string{"code": tc.code})
			assert.NoError(t, err)
			rec := httptest.NewRecorder()
			handler.Delete(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rec.Code)
			}

		})
	}
}

func TestHandlerGet(t *testing.T) {
	mockService := new(mocks.IURLService)

	type testCase struct {
		desc           string
		body, code     string
		mockSetup      func()
		expectedStatus int
	}

	testCases := []testCase{
		{
			desc: "success case",
			body: `{"url": "https://www.example.co"}`,
			code: "xyx",
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Get", mock.Anything).Return(&model.ShortURL{ID: 1},nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			desc: "failure case",
			body: `{"url": "https://www.example.co"}`,
			code: "xyx",
			mockSetup: func() {
				mockService = new(mocks.IURLService)

				mockService.On("Get", mock.Anything).Return(nil, errors.New("Please try again"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()
			h := URLHandler{service: mockService}
			handler := NewHandler(h.service)
			req, err := http.NewRequest(http.MethodDelete, "/shorten/"+tc.code, bytes.NewReader([]byte(tc.body)))
			req = muxSetVars(req, map[string]string{"code": tc.code})
			assert.NoError(t, err)
			rec := httptest.NewRecorder()
			handler.Get(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rec.Code)
			}

		})
	}
}