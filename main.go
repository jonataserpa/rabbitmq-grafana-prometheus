package main

import (
    "log"

    "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func main() {
    // Connect to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    // Declare a queue
    q, err := ch.QueueDeclare(
        "test_queue", // Queue name
        true,        // Durable
        false,        // Delete when unused
        false,        // Exclusive
        false,        // No-wait
        nil,          // Arguments
    )
    failOnError(err, "Failed to declare a queue")

    // Publish a message
    body := "Hello from Go!"
    err = ch.Publish(
        "",         // Exchange
        q.Name,     // Routing key
        false,      // Mandatory
        false,      // Immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
        })
    log.Printf(" [x] Sent %s", body)

    // Consume messages
    msgs, err := ch.Consume(
        q.Name, // Queue
        "",     // Consumer
        false,   // Auto-Ack
        false,  // Exclusive
        false,  // No Local
        false,  // No Wait
        nil,    // Args
    )
    failOnError(err, "Failed to register a consumer")

    // Loop to receive messages
    go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
        }
    }()

    // Keep the application running
    log.Println(" [*] Waiting for messages. To exit press CTRL+C")
    select {}
}
