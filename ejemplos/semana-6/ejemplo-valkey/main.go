package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type Container struct {
	Name   string
	CPU    float64
	Memory float64
	Status string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	rdb := redis.NewClient(&redis.Options{
		Addr:     "valkey:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("Error conectando a Valkey:", err)
	}
	log.Println("Conectado a Valkey")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Daemon iniciado. Generando datos cada 20 segundos...")

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			containers := generateRandomContainers()
			for _, c := range containers {
				containerKey := "container:" + c.Name + ":" + time.Now().Format("20060102150405")
				containerData := map[string]interface{}{
					"name":   c.Name,
					"cpu":    c.CPU,
					"memory": c.Memory,
					"status": c.Status,
					"time":   time.Now().UnixMilli(),
				}

				err := rdb.HSet(ctx, containerKey, containerData).Err()
				if err != nil {
					log.Println("Error insertando:", err)
				} else {
					rdb.Expire(ctx, containerKey, 10*time.Minute)
					log.Printf("Insertado contenedor random: %s (CPU: %.2f, MEM: %.2f, Estado: %s)\n",
						c.Name, c.CPU, c.Memory, c.Status)
				}
			}
		case <-sigs:
			log.Println("Daemon detenido.")
			break loop
		}
	}
}

func generateRandomContainers() []Container {
	names := []string{"nginx", "redis", "mysql", "golang-app", "nodejs-app", "python-app", "java-app", "ruby-app", "postgres", "mongodb", "ubuntu", "alpine"}
	statuses := []string{"running", "stopped", "paused", "restarting", "dead", "zombie"}

	n := rand.Intn(4) + 1
	var containers []Container
	for i := 0; i < n; i++ {
		containers = append(containers, Container{
			Name:   names[rand.Intn(len(names))],
			CPU:    rand.Float64() * 100,
			Memory: rand.Float64() * 512,
			Status: statuses[rand.Intn(len(statuses))],
		})
	}
	return containers
}
