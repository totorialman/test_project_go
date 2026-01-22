package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/totorialman/test_project_go/internal/config"
	courierHandler "github.com/totorialman/test_project_go/internal/handler/courier"
	courierRepository "github.com/totorialman/test_project_go/internal/repository/courier"
	courierUsecase "github.com/totorialman/test_project_go/internal/usecase/courier"
	deliveryHandler "github.com/totorialman/test_project_go/internal/handler/delivery"
	deliveryRepository "github.com/totorialman/test_project_go/internal/repository/delivery"
	deliveryUsecase "github.com/totorialman/test_project_go/internal/usecase/delivery"
)

func main() {
	log.SetOutput(os.Stdout)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	servPort := ":" + os.Getenv("PORT")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool := config.MustInitDB(ctx)
	defer dbPool.Close()

	courierRepo := courierRepository.NewRepository(dbPool)
	courierUsecase := courierUsecase.NewService(courierRepo)
	courierHandler := courierHandler.NewHandler(courierUsecase)

	deliveryRepo := deliveryRepository.NewRepository(dbPool)
	deliveryUsecase := deliveryUsecase.NewUsecase(courierRepo, deliveryRepo)
	deliveryHandler := deliveryHandler.NewHandler(deliveryUsecase)

	r := mux.NewRouter()
	r.HandleFunc("/ping", ping).Methods("GET")
	r.HandleFunc("/healthcheck", healthcheck).Methods("HEAD")
	r.HandleFunc("/courier", courierHandler.Create).Methods("POST")
	r.HandleFunc("/courier", courierHandler.Update).Methods("PUT")
	r.HandleFunc("/courier/{id}", courierHandler.GetById).Methods("GET")
	r.HandleFunc("/couriers", courierHandler.GetAll).Methods("GET")
	r.HandleFunc("/delivery/assign", deliveryHandler.Create).Methods("POST")

	server := &http.Server{
		Addr:         servPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server started on %s\n", servPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAnd Serve error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped gracefully")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
