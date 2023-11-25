package server

import (
	"encoding/json"
	"fmt"
	"go-web/internal/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.HelloWorldHandler)
	r.HandleFunc("/health", s.healthHandler)
	r.HandleFunc("/product", s.getProducts).Methods(http.MethodGet)
	r.HandleFunc("/product", s.HandleCreateProduct).Methods(http.MethodPost)

	return r
}

// func defaultHandleFunc(f apiFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if err := f(w, r); err != nil {
// 			err := JSONRep(w, http.StatusBadRequest, ApiError{Error: err.Error()})

// 			if err != nil {
// 				log.Fatalf("error handling JSON marshal. Err: %v", err)
// 			}

// 			return
// 		}
// 	}

// }

func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {
	prods := new([]database.Product)

	err := s.db.GetProducts(prods)

	if err != nil {
		JSONRep(w, http.StatusBadRequest, ApiError{Error: "invalid request"})
		return
	}

	err = JSONRep(w, http.StatusOK, &ApiResp{
		Data:    prods,
		Message: "success",
		Status:  http.StatusOK,
	})

	if err != nil {
		// handle error
		log.Printf("Error encoding JSON: %v", err)
	}

}

func (s *Server) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {

	p := &database.Product{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		JSONRep(w, http.StatusBadRequest, ApiError{Error: "invalid request"})
		return
	}

	if err := s.db.CreateProduct(p); err != nil {
		JSONRep(w, http.StatusInternalServerError, ApiError{Error: "internal server error"})
		return
	}

	JSONRep(w, http.StatusCreated, map[string]string{"message": "product created"})

}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())
	reqUrl := r.URL.Path
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}
	fmt.Print(string(jsonResp))
	fmt.Print(reqUrl)
	_, _ = w.Write(jsonResp)
}

func JSONRep(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
