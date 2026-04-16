package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Platillo struct {
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Descripcion string  `json:"descripcion"`
	Categoria   string  `json:"categoria"`
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	url := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir canal: %v", err)
	}
	defer ch.Close()

	queueName := "food_orders"
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al declarar cola: %v", err)
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al registrarse como consumer: %v", err)
	}

	log.Println("Consumer iniciado. Esperando mensajes...")

	for d := range msgs {
		var platillo Platillo
		err := json.Unmarshal(d.Body, &platillo)
		if err != nil {
			log.Printf("Error al unmarshal: %v", err)
			d.Nack(false, false)
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] ===== PEDIDO RECIBIDO =====\n", timestamp)
		fmt.Printf("  Platillo: %s\n", platillo.Nombre)
		fmt.Printf("  Precio: $%.2f\n", platillo.Precio)
		fmt.Printf("  Descripcion: %s\n", platillo.Descripcion)
		fmt.Printf("  Categoria: %s\n", platillo.Categoria)
		fmt.Printf("================================\n\n")

		d.Ack(false)
	}
}
