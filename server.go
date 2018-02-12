package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"./externalservice"
	"./utils"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Client externalservice.Client
}

func (s *Server) Run(addr string) {
	log.Println("Started API Server at", addr)
	if err := http.ListenAndServe(addr, utils.Logger()(s.Router)); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) InitializeRoutes() {
	s.Router = mux.NewRouter().StrictSlash(true)
	s.Router.HandleFunc("/api/posts/{id:[0-9]+}", s.get).Methods("GET")
	s.Router.HandleFunc("/api/posts/{id:[0-9]+}", s.post).Methods("POST")
}

func main() {
	s := Server{
		Client: &externalservice.ClientImpl{
			Posts: make(map[int]*externalservice.Post),
		},
	}
	s.InitializeRoutes()
	s.Run(":8080")
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	getPost, err := s.Client.GET(id)
	if err != nil {
		s.errorResponse(w, "Bad Request", fmt.Sprintf("/api/posts/%d", id), http.StatusBadRequest)
		return
	}
	s.successResponse(w, getPost, fmt.Sprintf("/api/posts/%d", id), http.StatusOK)
	return
}

func (s *Server) post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	title := r.FormValue("title")
	description := r.FormValue("description")
	postPayload := &externalservice.Post{
		ID:          id,
		Title:       title,
		Description: description,
	}
	savedPost, err := s.Client.POST(id, postPayload)
	if err != nil {
		s.errorResponse(w, err.Error(), fmt.Sprintf("/api/posts/%d", id), http.StatusBadRequest)
		return
	}
	s.successResponse(w, savedPost, fmt.Sprintf("/api/posts/%d", id), http.StatusCreated)
	return
}

func (s *Server) successResponse(w http.ResponseWriter, message interface{}, path string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    code,
		"message": message,
		"path":    path,
	})
}

func (s *Server) errorResponse(w http.ResponseWriter, message, path string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    code,
		"message": message,
		"path":    path,
	})
}
