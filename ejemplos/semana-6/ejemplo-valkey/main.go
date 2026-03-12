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

const TOTAL_RAM_GB = 6.4

type Container struct {
	Name   string
	CPU    float64
	Memory float64
	Status string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("Error conectando a Valkey:", err)
	}
	log.Println("Conectado a Valkey")

	// Guardar total RAM fijo una sola vez
	err := rdb.Set(ctx, "system:total_ram_gb", TOTAL_RAM_GB, 0).Err()
	if err != nil {
		log.Println("Error guardando total RAM:", err)
	} else {
		log.Printf("Total RAM guardada: %.1f GB\n", TOTAL_RAM_GB)
	}

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

			// ejecutar un comando que haga cd a esta carpeta: 'C:\Users\EQUIPO\Desktop\LAB-SO1-1S2026\ejemplos\semana-6\ejemplo-valkey\grafana-valkey' y luego ejecute 'docker compose up -d'
			//
			// cmd := exec.Command("docker", "compose", "up", "-d")
			// cmd.Dir = "C:\\Users\\EQUIPO\\Desktop\\LAB-SO1-1S2026\\ejemplos\\semana-6\\ejemplo-valkey\\grafana-valkey"
			// err := cmd.Run()
			// if err != nil {
			// 	log.Println("Error ejecutando docker compose:", err)
			// }

			// Guardar cada contenedor en su stream
			for _, c := range containers {
				err := rdb.XAdd(ctx, &redis.XAddArgs{
					Stream: "container:" + c.Name,
					MaxLen: 1000,
					Values: map[string]interface{}{
						"name":   c.Name,
						"cpu":    c.CPU,
						"memory": c.Memory,
						"status": c.Status,
					},
				}).Err()
				if err != nil {
					log.Println("Error insertando en stream:", err)
				} else {
					log.Printf("Insertado en stream: %s (CPU: %.2f, MEM: %.2f MB, Estado: %s)\n",
						c.Name, c.CPU, c.Memory, c.Status)
				}




				// Actualizar snapshot actual del contenedor (para top 5)
				rdb.HSet(ctx, "snapshot:"+c.Name, map[string]interface{}{
					"name":   c.Name,
					"cpu":    c.CPU,
					"memory": c.Memory,
					"status": c.Status,
				})

				// Sorted set para top 5 RAM
				rdb.ZAdd(ctx, "top:ram", redis.Z{
					Score:  c.Memory,
					Member: c.Name,
				})

				// Sorted set para top 5 CPU
				rdb.ZAdd(ctx, "top:cpu", redis.Z{
					Score:  c.CPU,
					Member: c.Name,
				})
			}

			// Calcular RAM usada total (suma de snapshots actuales)
			ramUsada := calcularRAMUsada(ctx, rdb)
			rdb.Set(ctx, "system:ram_used_mb", ramUsada, 0)
			ramLibre := (TOTAL_RAM_GB * 1024) - ramUsada
			if ramLibre < 0 {
				ramLibre = 0
			}
			rdb.Set(ctx, "system:ram_free_mb", ramLibre, 0)
			log.Printf("RAM usada: %.2f MB | RAM libre: %.2f MB\n", ramUsada, ramLibre)

		case <-sigs:
			log.Println("Daemon detenido.")
			break loop
		}
	}
}

func calcularRAMUsada(ctx context.Context, rdb *redis.Client) float64 {
	names := []string{"nginx", "redis", "mysql", "golang-app", "nodejs-app", "python-app", "java-app", "ruby-app", "postgres", "mongodb", "ubuntu", "alpine"}
	var total float64
	for _, name := range names {
		val, err := rdb.HGet(ctx, "snapshot:"+name, "memory").Float64()
		if err == nil {
			total += val
		}
	}
	return total
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
