package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"net/http"
	"os"
	_ "github.com/lib/pq"
)

var (
	queueTable = "message_queue"
	desiredMessageCount = 1
)

func main() {
	log.Println("Consumer started")

	http.HandleFunc("/health", func(http.ResponseWriter, *http.Request){})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	dbConnString := os.Getenv("DB_CONN_STRING")
	log.Println("dbConnString")
	log.Println(dbConnString)

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
	err := db.QueryRowContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = (SELECT id FROM %s ORDER BY timestamp LIMIT 1) RETURNING message", queueTable, queueTable)).Scan(&message)
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