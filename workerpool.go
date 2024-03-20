package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// task definition
type Task struct {
	InputData map[string]string
}

// way to process the task
func (t *Task) Process() {

	// map of input and output json keys
	var key_map = map[string]string{
		"ev":    "event",
		"et":    "event_type",
		"id":    "app_id",
		"uid":   "user_id",
		"mid":   "message_id",
		"t":     "page_title",
		"p":     "page_url",
		"l":     "browser_language",
		"sc":    "screen_size",
		"atrk":  "attribute_key",
		"atrv":  "attribute_value",
		"atrt":  "attribute_type",
		"uatrk": "user_trait_key",
		"uatrv": "user_trait_value",
		"uatrt": "user_trait_type",
	}
	// output data
	op_data := make(map[string]interface{})

	// maps to save list attributes and traits
	a_keys := make(map[int]string)
	a_values := make(map[int]string)
	a_types := make(map[int]string)
	ut_keys := make(map[int]string)
	ut_values := make(map[int]string)
	ut_types := make(map[int]string)

	for k, v := range t.InputData {
		if k[0] == 'a' {
			i, err := strconv.Atoi(k[4:])
			if err != nil {
				panic(err)
			}
			if k[3] == 'k' {
				a_keys[i] = v
			} else if k[3] == 'v' {
				a_values[i] = v
			} else if k[3] == 't' {
				a_types[i] = v
			}
		} else if k[0] == 'u' {
			if k[1] == 'a' {
				i, err := strconv.Atoi(k[5:])
				if err != nil {
					panic(err)
				}
				if k[4] == 'k' {
					ut_keys[i] = v
				} else if k[4] == 'v' {
					ut_values[i] = v
				} else if k[4] == 't' {
					ut_types[i] = v
				}
			}
		} else {
			op_data[key_map[k]] = v
		}
	}

	// iterate over atribute maps
	attr := make(map[string]map[string]string)
	for i := 1; i <= len(a_keys); i = i + 1 {
		x := make(map[string]string)
		x["value"] = a_values[i]
		x["type"] = a_types[i]
		attr[a_keys[i]] = x
	}

	//iterate over trait maps
	trt := make(map[string]map[string]string)
	for i := 1; i <= len(ut_keys); i = i + 1 {
		y := make(map[string]string)
		y["value"] = ut_values[i]
		y["type"] = ut_types[i]
		trt[ut_keys[i]] = y
	}

	op_data["attributes"] = attr
	op_data["traits"] = trt

	res, err := json.Marshal(op_data)
	if err != nil {
		log.Fatal(err)
	}
	// write formatted data

	resp, err := http.Post("https://webhook.site/7385d1e9-47cb-43bf-bbf7-79388cc095b6", "application/json", bytes.NewBuffer(res))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	//time.Sleep(2 * time.Second)
}

// worker pool definition
type WorkerPool struct {
	Tasks       []Task
	concurrency int
	tasksChan   chan Task
	wg          sync.WaitGroup

}

// Functions to execute the worker pool
func (wp *WorkerPool) worker() {

	for task := range wp.tasksChan {
		fmt.Printf("Processing the task...\n")
		task.Process()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) StartWorker() {
	//initialize the task channel
	wp.tasksChan = make(chan Task, 1)
	// start the worker go routines
	for i := 0; i < wp.concurrency; i++ {
		go wp.worker()
	}
}
func (wp *WorkerPool) Run() {

	wp.wg.Add(len(wp.Tasks))
	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}

	wp.wg.Wait()
}
