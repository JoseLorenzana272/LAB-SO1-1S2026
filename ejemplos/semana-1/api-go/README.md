# API Simple con Go y Fiber

Este ejemplo implementa una API REST básica utilizando el framework Fiber v2 en Go.

## Características
- Servidor web ligero con Fiber v2.52.10
- Middleware de logging para todas las peticiones
- Endpoint GET "/" que responde con "Hola, SOPES 1"
- Soporte para ejecución directa y en contenedor Docker

## Tecnologías
- Go 1.25.1
- Fiber v2 (framework web de alto rendimiento)
- Alpine Linux (para contenedores)

## Ejecución Directa
```bash
cd api-go
go run main.go
```

## Ejecución con Docker
```bash
cd api-go
docker build -t api-go:latest .
docker run -d -p 3000:3000 api-go:latest
```

## Pruebas
```bash
curl http://localhost:3000
```

El servidor responderá con: `Hola, SOPES 1`
