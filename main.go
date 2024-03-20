package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)



func GetData(w http.ResponseWriter, req *http.Request) {

	//Decode request body to map[string]string
	var ip_d map[string]string
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&ip_d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tasks []Task
	tasks = append(tasks, Task{InputData: ip_d})
	wp := WorkerPool{
		Tasks:       tasks,
		concurrency: 2,
	}
	wp.Run()
	fmt.Print("All tasks are completed\n")
}

func main() {

	http.HandleFunc("/read_input", GetData)

	fmt.Printf("Starting port server 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

