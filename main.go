package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Структура пристрою для Варіанта 4 [cite: 26]
type Device struct {
	Name        string    `json:"name"`
	DeviceType  string    `json:"device_type"`
	IPAddress   string    `json:"ip_address"`
	RoutingType string    `json:"routing_type"`
	Timestamp   time.Time `json:"timestamp"`
}

const logFileName = "network_log.txt"

func handleRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // [cite: 23]
		data, err := os.ReadFile(logFileName)
		if err != nil {
			fmt.Println("Лог порожній")
			w.Write([]byte("Лог-файл порожній"))
			return
		}
		fmt.Printf("\n--- Вміст логу ---\n%s", string(data))
		w.Write(data)

	case http.MethodPost: // [cite: 25, 27]
		var dev Device
		if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dev.Timestamp = time.Now()
		entry := fmt.Sprintf("[%s] %s (%s) IP:%s Rout:%s\n",
			dev.Timestamp.Format("15:04:05"), dev.Name, dev.DeviceType, dev.IPAddress, dev.RoutingType)

		f, _ := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(entry)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Запис додано")

	case http.MethodDelete: //
		os.Remove(logFileName)
		w.Write([]byte("Очищено"))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handleRequests)
	fmt.Println("Сервер запущено на http://localhost:8080/") // [cite: 22]
	log.Fatal(http.ListenAndServe(":8080", nil))
}
