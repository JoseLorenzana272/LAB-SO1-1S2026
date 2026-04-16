package main

import (
	"context"
	"encoding/json"
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

var platillos = []Platillo{
	{Nombre: "Tacos al Pastor", Precio: 45.00, Descripcion: "Tacos de cerdo adobado con tortilla de maize", Categoria: "Mexicana"},
	{Nombre: "Enchiladas Verdes", Precio: 85.00, Descripcion: "Enchiladas con salsa verde y queso", Categoria: "Mexicana"},
	{Nombre: "Burrito de Carne", Precio: 65.00, Descripcion: "Burrito grande con carne asada, frijoles y arroz", Categoria: "Tex-Mex"},
	{Nombre: "Guacamole con Totopos", Precio: 55.00, Descripcion: "Guacamole fresco con totopos caseros", Categoria: "Mexicana"},
	{Nombre: "Pozole Rojo", Precio: 95.00, Descripcion: "Pozole de cerdo conchile y rábano", Categoria: "Mexicana"},
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

	log.Println("Producer iniciado. Enviando platillos cada 25 segundos...")

	for i := 0; ; i++ {
		platillo := platillos[i%len(platillos)]

		body, err := json.Marshal(platillo)
		if err != nil {
			log.Printf("Error al marshal platillo: %v", err)
			continue
		}

		err = ch.PublishWithContext(context.Background(), "", queueName, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		if err != nil {
			log.Printf("Error al publicar: %v", err)
		} else {
			log.Printf("[PRODUCER] Enviado: %s - $%.2f", platillo.Nombre, platillo.Precio)
		}

		time.Sleep(25 * time.Second)
	}
}
