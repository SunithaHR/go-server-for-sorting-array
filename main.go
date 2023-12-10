// main.go
package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func sortSequential(arrays [][]int) [][]int {
	sortedArrays := make([][]int, len(arrays))
	for i, arr := range arrays {
		sorted := make([]int, len(arr))
		copy(sorted, arr)
		sort.Ints(sorted)
		sortedArrays[i] = sorted
	}
	return sortedArrays
}

func sortConcurrent(arrays [][]int) [][]int {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	sortedArrays := make([][]int, len(arrays))
	for i, arr := range arrays {
		wg.Add(1)
		go func(i int, arr []int) {
			defer wg.Done()
			sorted := make([]int, len(arr))
			copy(sorted, arr)
			sort.Ints(sorted)

			mutex.Lock()
			sortedArrays[i] = sorted
			mutex.Unlock()
		}(i, arr)
	}

	wg.Wait()
	return sortedArrays
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortSequential(requestPayload.ToSort)
	timeTaken := time.Since(startTime).Nanoseconds()

	responsePayload := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortConcurrent(requestPayload.ToSort)
	timeTaken := time.Since(startTime).Nanoseconds()

	responsePayload := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func main() {
	http.HandleFunc("/process-single", processSingleHandler)
	http.HandleFunc("/process-concurrent", processConcurrentHandler)

	http.ListenAndServe(":8000", nil)
}








