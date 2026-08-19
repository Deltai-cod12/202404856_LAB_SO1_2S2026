package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Confi de API1
const (
	carnet = "202404856"
	vm     = "VM1"
	port   = ":8081"
)

// HealthResponse define la estructura para el endpoint de salud
type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        string `json:"VM"`
	Carnet    string `json:"carnet"`
}

// ConnectionResponse define la estructura exacta solicitada por el proyecto para las llamadas cruzadas
type ConnectionResponse struct {
	APIName    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

// healthHandler maneja la peticion de estado de la API1
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "UP",
		Message:   "API1 is Ready",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		VM:        vm,
		Carnet:    carnet,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// callAPI2Handler realiza la llamada hacia la API2 ubicada en la VM1 (192.168.122.210)
func callAPI2Handler(w http.ResponseWriter, r *http.Request) {
	checkAPI(w, "API2", "VM1", "http://192.168.122.210:8082/health")
}

// callAPI3Handler realiza la llamada hacia la API3 ubicada en la VM2 (192.168.122.56)
func callAPI3Handler(w http.ResponseWriter, r *http.Request) {
	checkAPI(w, "API3", "VM2", "http://192.168.122.56:8083/health")
}

// checkAPI envia una peticion HTTP GET para validar si la API de destino esta activa
func checkAPI(w http.ResponseWriter, apiName string, targetVM string, url string) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get(url)

	if err != nil {
		sendConnectionResponse(w, apiName, targetVM, false)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var health HealthResponse
		err := json.NewDecoder(response.Body).Decode(&health)

		if err == nil && health.Status == "UP" {
			sendConnectionResponse(w, apiName, targetVM, true)
			return
		}
	}

	sendConnectionResponse(w, apiName, targetVM, false)
}

// sendConnectionResponse arma la respuesta JSON detallando la VM correspondiente segun el enunciado
func sendConnectionResponse(w http.ResponseWriter, apiName string, targetVM string, connected bool) {
	var message string

	if connected {
		message = fmt.Sprintf("The %s located on the %s is working", apiName, targetVM)
	} else {
		message = fmt.Sprintf("ERROR: The %s located on the %s is not working", apiName, targetVM)
	}

	response := ConnectionResponse{
		APIName:    apiName,
		Message:    message,
		Connection: connected,
		Carnet:     carnet,
	}

	w.Header().Set("Content-Type", "application/json")

	if !connected {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

// main configura los endpoints y levanto el servidor HTTP en el puerto indicado
func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api1/"+carnet+"/call-api2", callAPI2Handler)
	http.HandleFunc("/api1/"+carnet+"/call-api3", callAPI3Handler)

	log.Printf("API1 corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}