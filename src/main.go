package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

type Reminder struct {
	Text string `json:"text"`
}

var reminders []Reminder
var mu sync.Mutex

// Handler to set a new reminder
func setReminderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var rem Reminder
	if err := json.NewDecoder(r.Body).Decode(&rem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	reminders = append(reminders, rem)
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reminder set!"))
}

// Handler to get all reminders
func getRemindersHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default for local testing
	}

	http.HandleFunc("/set-reminder", setReminderHandler)
	http.HandleFunc("/get-reminder", getRemindersHandler)

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
