# Productor/Consumer con RabbitMQ

Ejemplo de sistema de mensajes usando Go y RabbitMQ.

## Local con Docker Compose

```bash
cd semana-11
docker compose up --build -d
```

- Productor: Envía platillos cada 25 segundos
- Consumer: Recibe y muestra los pedidos en cuanto los vea
- RabbitMQ UI: http://localhost:15672 (guest/guest)

### Ver mensajes acumulados en la cola

Para ver cómo se acumulan los mensajes en la UI:

```bash
# Detener consumer temporalmente
docker compose stop consumer

# Después de ~1 minuto, revisar la UI en la cola "food_orders"
# Verás "Ready" aumentando

# Reiniciar consumer
docker compose start consumer
```

## Kubernetes

### ConfigMap y Secrets

En `k8s/00-config.yaml` se definen:
- **ConfigMap**: Variables no sensibles (RABBITMQ_URL)
- **Secret**: Credenciales (también RABBITMQ_URL por seguridad)

Los pods los consumen con `envFrom`:
```yaml
envFrom:
  - secretRef:
      name: app-secret
```

Esto inyecta las variables automáticamente en los contenedores.

1. Construir y cargar imágenes:
```bash
docker build -t food-producer ./producer
docker build -t food-consumer ./consumer
```

2. Aplicar recursos:
```bash
kubectl apply -k k8s/
```

3. Acceder a RabbitMQ UI:
```bash
kubectl port-forward svc/rabbitmq 30072:15672
# http://localhost:30072 (guest/guest)
```

4. Ver logs:
```bash
kubectl logs -l app=producer
kubectl logs -l app=consumer
```

5. Eliminar:
```bash
kubectl delete -k k8s/
```
