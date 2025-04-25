package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// WorkQueue is a buffered channel that holds all the (unprocessed) work requests
var WorkQueue = make(chan WorkRequest, 100)

// WorkerQueue is a queue holding available workers
// a worker will be picked to handle work from this queue
var WorkerQueue chan chan WorkRequest

type WorkRequest struct {
	Nome  string
	Delay time.Duration
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

func parseDelay(d string) (time.Duration, error) {
	// Parse the delay.
	delay, err := time.ParseDuration(d)
	if err != nil {
		return delay, fmt.Errorf("bad delay value: %s", err.Error())
	}

	// Check to make sure the delay is anywhere from 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		return delay, errors.New("the delay must be between 1 and 10 seconds, inclusively")
	}

	return delay, nil
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Start starts the worker and listens for work (till Stop is called)
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue to signal we are available to pick up work
			w.WorkerQueue <- w.Work
			select {
			case task := <-w.Work:
				// Receive a work request.
				fmt.Printf("worker%d: Received work request, delaying for %f seconds\n", w.ID, task.Delay.Seconds())

				// time.Sleep(work.Delay)
				ticker := time.NewTicker(task.Delay)
				defer ticker.Stop()
				<-ticker.C
				// simulate actual work
				fmt.Printf("worker%d: Hello, %s!\n", w.ID, task.Nome)
				// can stop here on certain conditions
			case <-w.QuitChan:
				// check the "nil video by the Catalan guy, he had a scenario for a merge() where this didnt exactly work ok"
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop is used to notify the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// Collector simply handles and registers valid work requests over http
// we will ASSUME the clients issuing work requests NEVER wait for the work request to finish.
// This is an assumption to make our example simpler.
func Collector(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	delay, err := parseDelay(r.FormValue("delay"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	work := WorkRequest{Nome: name, Delay: delay}

	// Push the work onto the queue.
	WorkQueue <- work
	fmt.Println("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
}

func StartDispatcher(nWorkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nWorkers)
	limiter := make(chan struct{}, nWorkers*8)
	for i := 0; i < nWorkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		// This (ranging over WorkQueue) is buffered so we are ok, no need for external limiter
		for work := range WorkQueue {
			fmt.Println("Received work requeust")
			limiter <- struct{}{} // limit work execution to number of Workers
			go func(work WorkRequest) {
				defer func() { <-limiter }()
				worker := <-WorkerQueue

				fmt.Println("Dispatching work request")
				worker <- work
			}(work)
		}
	}()
}

var (
	NWorkers = flag.Int("numWorkers", runtime.NumCPU(), "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

func main() {
	flag.Parse()

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	StartDispatcher(*NWorkers)
	http.HandleFunc("/work", Collector)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}

// go run client_no_wait.go
// for i in {1..400};do curl localhost:8000/work -d name=$USER -d delay=$(expr $i % 20)s; done
