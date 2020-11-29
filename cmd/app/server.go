package app

import (
	"github.com/iamgafurov/crud/pkg/customers"
	_ "github.com/jackc/pgx/v4/stdlib"
	"strconv"
	"log"
	"net/http"
	"encoding/json"
	"errors"
)
 
type Server struct {
	mux *http.ServeMux
	customersSvc *customers.Service
}

func NewServer(mux *http.ServeMux, customersSvc *customers.Service) *Server{
	return &Server{mux:mux, customersSvc:customersSvc}
}

func (s *Server)ServeHTTP(writer http.ResponseWriter, request *http.Request){
	s.mux.ServeHTTP(writer,request)
}

func (s *Server) Init(){
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomersByID)
	s.mux.HandleFunc("/customers.getAll", s.handleGetCustomersAll)
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetCustomersAllActive)
	s.mux.HandleFunc("/customers.save", s.handleGetCustomersSave)
	s.mux.HandleFunc("/customers.removeById", s.handleCustomersRemoveByID)
	s.mux.HandleFunc("/customers.blockById", s.handleCustomersBlockByID)
	s.mux.HandleFunc("/customers.unblockById", s.handleCustomersUnBlockByID)
}

func (s *Server) handleGetCustomersByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Print(err)
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	item,err := s.customersSvc.ByID(request.Context(), id)
	if errors.Is(err,customers.ErrNotFound){
		http.Error(writer, http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}

	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	writer.Header().Set("Content-Type", "application/json")
	_,err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetCustomersAll(w http.ResponseWriter, r *http.Request){
	items, err := s.customersSvc.All(r.Context())
	if errors.Is(err, customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound),http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}

	data,err := json.Marshal(items)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
}

func (s *Server) handleGetCustomersAllActive(w http.ResponseWriter, r *http.Request){
	items, err := s.customersSvc.AllActive(r.Context())
	if errors.Is(err, customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound),http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}

	data,err := json.Marshal(items)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
}

func (s *Server) handleGetCustomersSave(w http.ResponseWriter, r *http.Request){
	var customer *customers.Customer
	idParam := r.FormValue("id")
	log.Print(idParam)
	nameParam := r.FormValue("name")
	log.Print(nameParam)
	phoneParam := r.FormValue("phone")
	log.Print(phoneParam)
	id,err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	if id == 0 {
		customer, err= s.customersSvc.Create(r.Context(),nameParam,phoneParam)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
			return 
		}
	}else {
		customer, err= s.customersSvc.Update(r.Context(),id,nameParam,phoneParam)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
			return 
		}
	}

	data,err := json.Marshal(customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
}

func (s *Server) handleCustomersRemoveByID(w http.ResponseWriter, r *http.Request){
	idParam:= r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Print(err)
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	customer,err := s.customersSvc.RemoveByID(r.Context(),id)
	if errors.Is(err,customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	data,err := json.Marshal(customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
}


func (s *Server) handleCustomersBlockByID(w http.ResponseWriter, r *http.Request){
	idParam:= r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Print(err)
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	customer,err := s.customersSvc.BlockByID(r.Context(),id)
	if errors.Is(err,customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	data,err := json.Marshal(customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}

}

func (s *Server) handleCustomersUnBlockByID(w http.ResponseWriter, r *http.Request){
	idParam:= r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Print(err)
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	customer,err := s.customersSvc.UnBlockByID(r.Context(),id)
	if errors.Is(err,customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	data,err := json.Marshal(customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return 
	}

}