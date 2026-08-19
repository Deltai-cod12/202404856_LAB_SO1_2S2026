package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// COnfi de api2 y las rutas
const (
    carnet = "202404856"
    vm     = "VM1"
    port   = ":8082"

    api1URL = "http://192.168.122.210:8081"
    api3URL = "http://192.168.122.56:8083"
)

// Estructura de health
type HealthResponse struct {
    Status    string `json:"status"`
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
    VM        string `json:"VM"`
    Carnet    string `json:"carnet"`
}

// Estructura que uso cuando intento comunicarme con otra API
type ConnectionResponse struct {
    APIName    string `json:"apiname"`
    Message    string `json:"message"`
    Connection bool   `json:"connection"`
    Carnet     string `json:"carnet"`
}

// Endpoint de salud
func healthHandler(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{
        Status:    "UP",
        Message:   "API2 is Ready",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        VM:        vm,
        Carnet:    carnet,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Manejador para llamar a la API1
func callAPI1Handler(w http.ResponseWriter, r *http.Request) {
    checkAPI(w, "API1", "VM1", api1URL+"/health")
}

// Manejador para llamar a la API3
func callAPI3Handler(w http.ResponseWriter, r *http.Request) {
    checkAPI(w, "API3", "VM2", api3URL+"/health")
}

// Funcion generica donde valido si la API destino responde o no
func checkAPI(w http.ResponseWriter, apiName string, targetVM string, url string) {
    client := http.Client{
        Timeout: 5 * time.Second,
    }

    response, err := client.Get(url)

    // Si fallo la peticion, asumo que no hay conexion
    if err != nil {
        sendConnectionResponse(w, apiName, targetVM, false)
        return
    }

    defer response.Body.Close()

    // Si la respuesta fue exitosa (200 OK) y el estado es UP, marco conexion exitosa
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

// Construyo y envio la respuesta JSON indicando el estado 
func sendConnectionResponse(w http.ResponseWriter, apiName string, targetVM string, connected bool) {
    var message string

    if connected {
        message = fmt.Sprintf(
            "The %s located on the %s is working",
            apiName,
            targetVM,
        )
    } else {
        message = fmt.Sprintf(
            "ERROR: The %s located on the %s is not working",
            apiName,
            targetVM,
        )
    }

    response := ConnectionResponse{
        APIName:    apiName,
        Message:    message,
        Connection: connected,
        Carnet:     carnet,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // Configuro mis rutas de endpoints
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc(
        "/api2/"+carnet+"/call-api1",
        callAPI1Handler,
    )
    http.HandleFunc(
        "/api2/"+carnet+"/call-api3",
        callAPI3Handler,
    )
    log.Printf("API2 running on %s", port)
    // Arranco mi servidor en el puerto 8082
    log.Fatal(http.ListenAndServe(port, nil))
}