package router

import (
	"github.com/gorilla/mux"
	"robin-procedure-service-go/middleware"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/procedures", middleware.GetAllProcedures).Methods("GET", "OPTIONS")
	router.HandleFunc("/procedures/{id}", middleware.GetAllProcedures).Methods("GET", "OPTIONS")

	return router
}
