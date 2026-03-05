#!/bin/bash

# Detener el daemon de valkey si está corriendo
if ! sudo systemctl stop containers-daemon 2>/dev/null; then
    echo "[INFO] El daemon no estaba corriendo."
fi

if ! sudo systemctl disable containers-daemon 2>/dev/null; then
    echo "[INFO] El daemon no estaba habilitado."
fi

# ejecutar docker-compose down en la carpeta grafana-valkey
echo "--------------------------------"
echo "Deteniendo contenedores de docker-compose..."
if ! (cd ./grafana-valkey && docker-compose down); then
    echo "[ERROR] No se pudo detener los contenedores de docker-compose."
fi
echo "Contenedores detenidos."

# Build the Go daemon
echo "[INFO] Construyendo Go daemon..."
if ! go build -o daemon main.go; then
    echo "[ERROR] Falló la construcción del Go daemon. Abortando script."
    exit 1
fi

echo "[INFO] Proceso completado exitosamente."

echo "--------------------------------"

echo "Cargando servicio en systemd..."

if ! sudo systemctl daemon-reload; then
    echo "[ERROR] Falló la recarga de systemd."
    exit 1
fi

if ! sudo systemctl enable containers-daemon; then
    echo "[ERROR] Falló al habilitar el servicio containers-daemon."
    exit 1
fi

if ! sudo systemctl start containers-daemon; then
    echo "[ERROR] Falló al iniciar el servicio containers-daemon."
    exit 1
fi

echo "Servicio systemd creado e iniciado."

echo "--------------------------------"

echo "Iniciando contenedores de docker-compose..."
if ! (cd ./grafana-valkey && docker-compose up -d); then
    echo "[ERROR] No se pudo iniciar los contenedores de docker-compose."
    exit 1
fi
echo "Contenedores iniciados."

echo "--------------------------------"
echo "Script completado exitosamente."
