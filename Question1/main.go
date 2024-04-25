package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	numbers []int
	window  = 10
)
func handleNumbers(w http.ResponseWriter, r *http.Request) {
	numberID := r.URL.Path[len("/numbers/"):]
	if !isValidNumberID(numberID) {
		http.Error(w, "Invalid number ID", http.StatusBadRequest)
		return
	}

	nums, err := fetchNumbers(numberID)
	if err != nil {
		http.Error(w, "Error fetching numbers from test server", http.StatusInternalServerError)
		return
	}

	updateWindow(nums)

	resp := map[string]interface{}{
		"numbers":          nums,
		"windowPrevState":  numbers[:len(numbers)-len(nums)],
		"windowCurrState":  numbers,
		"avg":              calculateAverage(numbers),
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func isValidNumberID(numberID string) bool {
	switch numberID {
	case "p", "f", "e", "r":
		return true
	default:
		return false
	}
}

func fetchNumbers(numberID string) ([]int, error) {
	
	
	switch numberID {
	case "p": // Prime numbers
		return []int{2, 3, 5, 7, 11}, nil
	case "f": // Fibonacci numbers
		return []int{0, 1, 1, 2, 3}, nil
	case "e": // Even numbers
		return []int{2, 4, 6, 8, 10}, nil
	case "r": // Random numbers
		return []int{7, 3, 9, 5, 1}, nil
	default:
		return nil, fmt.Errorf("unsupported number ID")
	}
}

func updateWindow(newNumbers []int) {
	mu.Lock()
	defer mu.Unlock()

	numbers = append(numbers, newNumbers...)
	if len(numbers) > window {
		numbers = numbers[len(numbers)-window:]
	}
}

func calculateAverage(nums []int) float64 {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return float64(sum) / float64(len(nums))
}


func main() {
	router := http.NewServeMux()
	router.HandleFunc("/numbers/", handleNumbers)

	server := &http.Server{
		Addr:         "localhost:9876",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

