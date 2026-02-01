# Fiber API Go

API REST desarrollada con Go y el framework Fiber para la asignatura de Sistemas Operativos 1.

## Características

- API REST ligera y de alto rendimiento
- Endpoint de health check
- Middleware de logging HTTP
- Soporte para Docker con multi-stage build
- Imagen final optimizada con Alpine Linux

## Requisitos Previos

- Go 1.25.1 o superior
- Docker (opcional, para ejecución en contenedor)

## Instalación

1. Clonar el repositorio
2. Navegar al directorio del proyecto:
```bash
cd fiber-api-go
```

3. Instalar dependencias:
```bash
go mod download
```

## Ejecución Local

Ejecutar la aplicación directamente con Go:

```bash
go run main.go
```

La API estará disponible en `http://localhost:8081`

## Ejecución con Docker

Construir la imagen Docker:

```bash
docker build -t fiber-api .
```

Ejecutar el contenedor:

```bash
docker run -p 8081:8081 fiber-api
```

## Endpoints de la API

### Health Check
**GET** `/health`

Verifica el estado de la API.

**Respuesta:**
```json
{
  "status": "UP",
  "message": "API4 is Ready"
}
```

**Curl:**
```bash
curl http://localhost:8081/health
```

## Estructura del Proyecto

```
fiber-api-go/
├── main.go          # Punto de entrada de la aplicación
├── go.mod           # Definición de dependencias
├── go.sum           # Checksums de dependencias
├── dockerfile       # Configuración de Docker
└── .dockerignore    # Archivos a ignorar en Docker
```

## Stack Tecnológico

- **Go 1.25.1**: Lenguaje de programación principal
- **Fiber v2**: Framework web de alto rendimiento
- **Alpine Linux 3.22**: Imagen base ligera para Docker
- **Docker**: Contenedorización de la aplicación

## Detalles de la Imagen Docker

La imagen Docker utiliza un multi-stage build:

1. **Stage builder**: Imagen `golang:1.25-alpine` para compilar la aplicación
2. **Stage final**: Imagen `alpine:3.22` para ejecutar el binario compilado

El binario compilado ocupa significativamente menos espacio y solo contiene lo necesario para la ejecución.
