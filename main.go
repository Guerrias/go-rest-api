package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guerrias/go-rest-api/employee"
)

func main() {
	// init router
	r := mux.NewRouter()

	//routes
	r.HandleFunc("/api/employees", employee.GetEmployees).Methods("Get")
	r.HandleFunc("/api/employees/{id}", employee.GetEmployee).Methods("Get")
	r.HandleFunc("/api/employees", employee.CreateEmployee).Methods("Post")
	r.HandleFunc("/api/employees/{id}", employee.UpdateEmployee).Methods("Put")
	r.HandleFunc("/api/employees/{id}", employee.DeleteEmployee).Methods("Delete")

	log.Fatal(http.ListenAndServe(":8000", r))
}
