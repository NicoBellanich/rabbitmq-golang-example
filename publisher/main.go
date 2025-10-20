package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	if rabbitURL == "" {
		rabbitURL = "amqp://user:pass@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v - rabbitURL (%v)", err, rabbitURL)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("tasks", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		body := r.URL.Query().Get("msg")
		if body == "" {
			http.Error(w, "missing msg parameter", 400)
			return
		}
		err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
		if err != nil {
			http.Error(w, "Failed to publish", 500)
			return
		}
		fmt.Fprintf(w, "Sent: %s", body)
	})

	log.Println("Publisher running on :8080")
	http.ListenAndServe(":8080", nil)
}
