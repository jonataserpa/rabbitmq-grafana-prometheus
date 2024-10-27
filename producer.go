package main

import (
    "log"
    // "time"
    "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "orders", true, false, false, false, nil,
    )
    failOnError(err, "Failed to declare a queue")

    for i := 0; i < 1000000; i++ {
        //time.Sleep(500 * time.Millisecond)
        body := "Teste de public, RabbitMQ "
        err = ch.Publish(
            "", q.Name, false, false,
            amqp.Publishing{ContentType: "text/plain", Body: []byte(body)},
        )
        failOnError(err, "Failed to publish a message")
        log.Printf(" [x] Sent %s", body)
    }
}
