package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"math/rand"
	"strings"

	_ "github.com/lib/pq"
)

// "os/exec"
// "bytes"
// "regexp"

func main() {
	http.HandleFunc("/enqueue", enqueueMessageHandler)
	http.HandleFunc("/health", func(http.ResponseWriter, *http.Request){})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type RequestBody struct {
	Message string `json:"message"`
}

var (
	maxAgents = 10
)

func enqueueMessageHandler(w http.ResponseWriter, r *http.Request) {
	
	dbConnString := os.Getenv("DB_CONN_STRING")
	queueTable := os.Getenv("QUEUE_TABLE")

	// An example of the connection string we expect
	//  postgresql://consumer:pass@keda-postgres.keda-demo.svc.cluster.local:80/queue?sslmode=disable
	parts := strings.Split(dbConnString, ":")
	workerId := rand.Intn(maxAgents) + 1

	// Insert a worker id to better parallelize production. SQL only allows sequential logins by the same user
	dbConnString = fmt.Sprintf("%s:%s%d:%s:%s", parts[0], parts[1], workerId, parts[2], parts[3])

	db, err := sql.Open("postgres", dbConnString)
	log.Println("DB Connection string: ", dbConnString)
	log.Println("DB Queue Table: ", queueTable)
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestBody RequestBody
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Println("Error parsing request body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := requestBody.Message

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")


	_, err = db.ExecContext(r.Context(), fmt.Sprintf("INSERT INTO %s (message) VALUES ($1)", queueTable), message)
	if err != nil {
		log.Println("Error enqueueing message:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		_, err = db.ExecContext(r.Context(), fmt.Sprintf("SELECT COUNT(*) FROM %s", queueTable))
		return
	}

	log.Println("Message enqueued successfully")
	fmt.Fprintln(w, "Message enqueued successfully")
}