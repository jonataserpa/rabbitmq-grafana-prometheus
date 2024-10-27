package main

import (
    "database/sql"
    "log"
    // "time"
    "github.com/streadway/amqp"
    _ "github.com/go-sql-driver/mysql"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

// Consumer - consome mensagens e salva no banco de dados
func consumeMessages(ch *amqp.Channel, queueName string, db *sql.DB) {
    msgs, err := ch.Consume(
        "test_queue", 
        "", 
        false, 
        false, 
        false, 
        false, 
        nil,
    )
    failOnError(err, "Failed to register a consumer")

    // Número de consumidores (goroutines)
	numConsumers := 5
	for i := 0; i < numConsumers; i++ {
		go func(consumerID int) {
			for d := range msgs {
				log.Printf("Consumer %d: Received a message: %s", consumerID, d.Body)

				// Simular o processamento da mensagem (ex.: salvar no banco de dados)
				//time.Sleep(500 * time.Millisecond) // Simula processamento

                 // Insere a mensagem no banco de dados
                _, err := db.Exec("INSERT INTO orders (message) VALUES (?)", string(d.Body))
                if err != nil {
                    log.Printf("Failed to insert message: %v", err)
                    // Caso o salvamento falhe, rejeitamos a mensagem para reprocessamento
                    d.Nack(false, true) // Reinsere a mensagem na fila
                    continue
                }

                // Confirma manualmente a mensagem após salvar com sucesso no banco
                err = d.Ack(true)
                if err != nil {
                    log.Printf("Failed to ack message: %v", err)
                }
			}
		}(i + 1)
	}

	// Manter o programa em execução
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

func main() {
    // Configuração do RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")

    // Configura o prefetch para o canal
    err = ch.Qos(1, 0, false) // Prefetch de 10 mensagens
    failOnError(err, "Failed to set QoS")

    defer ch.Close()

    queue, err := ch.QueueDeclare(
        "test_queue", 
        true, 
        false, 
        false, 
        false, 
        nil,
    )
    failOnError(err, "Failed to declare a queue")

    // Configuração do MySQL
    db, err := sql.Open("mysql", "root:root@tcp(localhost:3308)/orders")
    failOnError(err, "Failed to connect to MySQL")
    defer db.Close()

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id INT AUTO_INCREMENT,
            message TEXT CHARACTER SET utf8mb4,
            PRIMARY KEY(id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `)
    failOnError(err, "Failed to create table in MySQL")

    // Inicia o producer e o consumer
    consumeMessages(ch, queue.Name, db)
}
