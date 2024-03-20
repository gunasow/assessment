package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)


var wp= WorkerPool{
	concurrency: 3,
}
func GetData(w http.ResponseWriter, req *http.Request) {

	//Decode request body to map[string]string
	var ip_d map[string]string
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&ip_d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wp.Tasks=[]Task{
		{
			InputData: ip_d,
		},
	}
	wp.Run()
	fmt.Print("Task completed\n")
	
}


func main() {

	http.HandleFunc("/read_input", GetData)

	//start worker
	wp.StartWorker()
	fmt.Printf("No of workers: %d\n",runtime.NumGoroutine())
	
	fmt.Printf("Starting port server 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

