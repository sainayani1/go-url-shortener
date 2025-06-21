package main

import (
	"log"
	"net/http"
	"url-shortener/app"
	"url-shortener/shortenurl/repo"
	"url-shortener/shortenurl/service"
	"url-shortener/shortenurl/handler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (loading from OS env)")
	}

	app.InitDB()

	repoLayer := repo.NewURLRepo(app.DB)
	svc := service.NewURLService(&repoLayer)
	h := handler.NewHandler(svc)
	
	r := mux.NewRouter()
	r.HandleFunc("/shorten", h.Create).Methods("POST")
	r.HandleFunc("/shorten/{code}", h.Get).Methods("GET")
	r.HandleFunc("/shorten/{code}", h.Update).Methods("PUT")
	r.HandleFunc("/shorten/{code}", h.Delete).Methods("DELETE")
	// r.HandleFunc("/shorten/{code}/stats", handlers.GetStats).Methods("GET")

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
