package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"net/http"
	"os"
	"strings"
	"math/rand"
	"strconv"
	"io/ioutil"
	_ "github.com/lib/pq"
)

var (
	queueTable = "message_queue"
	desiredMessageCount = 500
	maxAgents = 100
)

func getIndex(file) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	numStr =: strings.TrimSpace(b)
	i, err := strconv.Atoi(numStr)
	return i
}

func main() {
	log.Println("Consumer started")

	http.HandleFunc("/health", func(http.ResponseWriter, *http.Request){})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	dbConnString := os.Getenv("DB_CONN_STRING")
	log.Println("dbConnString")
	log.Println(dbConnString)

	// An example of the connection string we expect
	//  postgresql://consumer:pass@keda-postgres.keda-demo.svc.cluster.local:80/queue?sslmode=disable
	parts := strings.Split(dbConnString, ":")
	workerId := -1
	if (len(os.Args) > 0) {
		log.Println("Worker id was set already")
		workerId = getIndex(os.Args[0])
	} else {
		log.Println("Setting worker id randomly!")
		workerId = rand.Intn(maxAgents) + 1
	}

	// Insert a worker id to better parallelize consumption. SQL only allows sequential logins by the same user
	dbConnString = fmt.Sprintf("%s:%s%d:%s:%s", parts[0], parts[1], workerId, parts[2], parts[3])

	ctx := context.Background()

	log.Println("Opening postgres connection")

	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	
	log.Println("Listening for messages in the queue...")

	processedMessages := 0

	// This will only process a single message
	for processedMessages < desiredMessageCount {
		message, err := dequeueMessage(ctx, db)

		if err != nil {
			if err != sql.ErrNoRows {
				log.Println("Error dequeuing message:", err)
				os.Exit(1)
			} else {
				log.Println("No more messages to process, exiting...")
				return
			}

		} else {
			log.Println("Processing message:", message)
			processedMessages++
			// Simulate message processing
			if processMessageErr := processMessage(message); processMessageErr != nil {
				log.Println("Error processing message:", processMessageErr)
			} else {
				log.Println("Message processed successfully")
			}
		}
	}
	
	log.Println("Consumer finished")
}

func dequeueMessage(ctx context.Context, db *sql.DB) (string, error) {
	var message string
	log.Println("Dequeue message!")
	err := db.QueryRowContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = (SELECT id FROM %s ORDER BY RANDOM() LIMIT 1) RETURNING message", queueTable, queueTable)).Scan(&message)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("Error dequeuing message:", err)
		}
		return "", err
	}

	return message, nil
}

func processMessage(message string) error {
	// Simulated message processing
	time.Sleep(2 * time.Second)
	// Simulated processing error
	if message == "error" {
		return fmt.Errorf("simulated processing error")
	}
	return nil
}