package main

import (
	"cryptodashboard/internal/handlers"
	"cryptodashboard/internal/pubsub"
	"cryptodashboard/internal/services"
	"cryptodashboard/internal/worker"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

var (
	apiKey    = getEnvOrDefault("BUD_API_KEY", "una api key bonita")
	apiSecret = getEnvOrDefault("BUD_API_SECRET", "una api secret bonita")
	baseURL   = getEnvOrDefault("BUD_BASE_URL", "https://api.buda.com/v2")
)

func main() {
	// Configurar el servidor
	mux := http.NewServeMux()

	fmt.Println("apiKey", apiKey, "apiSecret", apiSecret, "baseURL", baseURL)
	// Servir archivos estáticos
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Crear instancias de servicios
	budaService := services.NewBuda(baseURL, apiKey, apiSecret)
	pubsubInstance := pubsub.NewPubSub()

	// Crear y configurar el worker
	budaWorker := worker.NewBudaWorker(budaService, pubsubInstance, 30*time.Second)
	budaWorker.Start()

	loc, _ := time.LoadLocation("America/Santiago")

	// Crear contexto para los handlers
	c := handlers.Context{
		BudaServices: budaService,
		PubSub:       pubsubInstance,
		Loc:          loc,
	}

	// Configurar rutas
	mux.HandleFunc("/dashboard", c.DashboardHandler)
	mux.HandleFunc("/balances", c.BalanceHandler)

	// Iniciar el servidor
	log.Println("Servidor iniciando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
