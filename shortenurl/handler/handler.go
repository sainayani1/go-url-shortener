package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"url-shortener/shortenurl/model"
	"url-shortener/shortenurl/service"

	"github.com/gorilla/mux"
)

type URLHandler struct {
	service service.IURLService
}

func NewHandler(s service.IURLService) *URLHandler {
	return &URLHandler{service: s}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := model.ShortURLCreateReq{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	url, err := h.service.Create(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, url)
}

func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	url,err := h.service.Get(code)
	switch {
	case err == nil:
		respondJSON(w, http.StatusOK, url)

	case strings.HasPrefix(err.Error(), "short URL not found"):
		http.Error(w, "Short URL not found", http.StatusNotFound)

	case strings.HasPrefix(err.Error(), "access count update failed"):
		http.Error(w, "Access count could not be updated, please retry", http.StatusInternalServerError)

	default:
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func(h *URLHandler) Update(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	req := model.ShortURLCreateReq{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	url, err := h.service.Update(code, req.URL)
	if err != nil {
		http.Error(w, "update failed "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, url)
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	err := h.service.Delete(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
