package router

import (
	"github.com/gorilla/mux"
	"robin-procedure-service-go/middleware"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/procedures", middleware.GetAllProcedures).Methods("GET", "OPTIONS")
	router.HandleFunc("/procedures/{id}", middleware.GetProcedure).Methods("GET", "OPTIONS")
	router.HandleFunc("/procedures", middleware.CreateProcedure).Methods("POST", "OPTIONS")
	router.HandleFunc("/procedures/{id}", middleware.UpdateProcedure).Methods("PUT", "OPTIONS")
	router.HandleFunc("/procedures/{id}", middleware.DeleteProcedure).Methods("DELETE", "OPTIONS")

	return router
}
