package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noonyuu/comparison/backend/internal/db"
)

type MovieHandler struct {
	DB *db.Database
}

func NewRouter(database *db.Database) *mux.Router {
	handler := &MovieHandler{DB: database}
	r := mux.NewRouter()
	r.HandleFunc("/api/movies", handler.GetMovies).Methods("GET")
	r.HandleFunc("/api/movies/{id}", handler.GetMovie).Methods("GET")
	r.HandleFunc("/api/movies", handler.CreateMovie).Methods("POST")
	return r
}

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.DB.GetMovies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	movie, err := h.DB.GetMovie(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie db.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.DB.CreateMovie(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}
