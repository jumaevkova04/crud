package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jumaevkova04/crud/cmd/app/middleware"
	"github.com/jumaevkova04/crud/pkg/customers"
	"github.com/jumaevkova04/crud/pkg/security"
)

// Server ...
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	securitySvc  *security.Service
}

// NewServer ...
func NewServer(mux *mux.Router, customersSvc *customers.Service, securitySvc *security.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc, securitySvc: securitySvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

func (s *Server) Init() {
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Basic(s.securitySvc.Auth))
	s.mux.Use(middleware.CheckHeader("Content-Type", "application/json"))
	// s.mux.Handle("/customers", middleware.Logger(http.HandlerFunc(s.handleGetAllCustomers))).Methods(GET)

	// chMd := middleware.CheckHeader("Content-Type", "application/json")
	// s.mux.Handle("/customers", chMd(http.HandlerFunc(s.handleGetAllCustomers))).Methods(GET)

	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers", s.handleSaveCustomer).Methods(POST)
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnblockByID).Methods(DELETE)
}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.ByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, item)
	if err != nil {
		log.Println("ERROR", err)
	}
}

func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customersSvc.All(r.Context())

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, items)
	if err != nil {
		log.Println("ERROR", err)
	}
}

func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customersSvc.AllActive(r.Context())
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, items)
	if err != nil {
		log.Println("ERROR", err)
	}
}

func (s *Server) handleSaveCustomer(w http.ResponseWriter, r *http.Request) {
	var customer *customers.Customer

	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.Save(r.Context(), customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, item)
	if err != nil {
		log.Println("ERROR", err)
	}
}

func (s *Server) handleRemoveByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.RemoveByID(r.Context(), id)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, item)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.BlockByID(r.Context(), id)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, item)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
}

func (s *Server) handleUnblockByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.UnblockByID(r.Context(), id)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, item)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
}

func sendResponse(w http.ResponseWriter, response interface{}) error {
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println("ERROR", err)
		return err
	}
	return nil
}
