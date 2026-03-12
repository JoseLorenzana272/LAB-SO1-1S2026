Este ejemplo hace lo siguiente:

1. Se conecta a una base de datos **Valkey** (Redis in-memory).
2. Cada **20 segundos**, genera hasta **4 contenedores con datos aleatorios** (nombre, CPU, memoria, estado) y los guarda como hashes en Valkey.
3. Se ejecuta como **daemon** en segundo plano.

---

## 📌 Compilación

En la carpeta del código:

```bash
go mod init daemon-test
go get github.com/redis/go-redis/v9
go build -o daemon main.go
```

---

## 📌 Ejecución

Asegúrate de tener Valkey corriendo:

```bash
docker run -d --name valkey -p 6379:6379 valkey/valkey:latest
```

Luego ejecutar el daemon:

```bash
./daemon &
```

Ver logs en tiempo real:

```bash
tail -f nohup.out
```

Detenerlo:

```bash
pkill daemon
```

---

## 📌 Instalar como Daemon en Linux (systemd)

Crea el archivo de servicio:

```bash
sudo nano /etc/systemd/system/grafana-db-daemon.service
```

Contenido:

```txt
[Unit]
Description=Daemon en Go - Valkey
After=network.target

[Service]
ExecStart=/ruta/a/tu/daemon
WorkingDirectory=/ruta/a/tu/carpeta
Restart=always

[Install]
WantedBy=multi-user.target
```

Guardar y habilitar:

```bash
sudo systemctl daemon-reload
sudo systemctl enable grafana-db-daemon
sudo systemctl start grafana-db-daemon
```

Ver estado:

```bash
sudo systemctl status grafana-db-daemon
```

Parar y quitar el daemon:
```bash
sudo systemctl stop grafana-db-daemon
sudo systemctl disable grafana-db-daemon
sudo rm /etc/systemd/system/grafana-db-daemon.service
sudo systemctl daemon-reload
```

---

## GRAFANA

### 📂 Carpeta grafana-valkey

En la carpeta `grafana-valkey` se encuentra la configuración de Grafana con el plugin Redis.

---

## 📌 Dockerfile de Grafana con Redis plugin

Grafana soporta **Redis/Valkey** como datasource mediante el plugin oficial `redis-datasource`.

```dockerfile
FROM grafana/grafana:11.0.0

RUN grafana-cli plugins install redis-datasource

VOLUME ["/var/lib/grafana"]

EXPOSE 3000
```

---

## 📌 docker-compose.yml

```yaml
services:
  valkey:
    image: valkey/valkey:latest
    container_name: valkey-db
    ports:
      - "6379:6379"
    volumes:
      - ./valkey-data:/data

  grafana:
    build: .
    container_name: grafana-valkey
    ports:
      - "3001:3000"
    user: "0"
    volumes:
      - ./grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_DATASOURCES_DEFAULT_NAME=valkey
    depends_on:
      - valkey
```

---

## 📌 Levantar Grafana

En la carpeta grafana-valkey:

```bash
docker-compose up -d
```

Grafana estará en:
*http://localhost:3001*
Usuario: `admin`
Password: `admin`

---

## 📌 Configurar Data Source Redis

1. En Grafana, ve a **Connections → Data Sources → Add new data source**.
2. Busca **Redis**.
3. Configura:
   - URL: `valkey:6379`
   - Typo: `Redis`
4. Guarda y prueba.

---

## 📌 Crear Dashboard

Puedes usar queries como:

```json
HGETALL container:*
```

O en modo Command:
```
HGETALL container:nginx:*


```

---

## Qué se guarda en Valkey ahora:

| Key | Tipo | Contenido |
|---|---|---|
| `system:total_ram_gb` | STRING | `6.4` fijo |
| `system:ram_used_mb` | STRING | suma RAM activa |
| `system:ram_free_mb` | STRING | RAM libre calculada |
| `container:nginx` | STREAM | historial CPU/RAM |
| `snapshot:nginx` | HASH | último valor de nginx |
| `top:ram` | ZSET | ranking por memoria |
| `top:cpu` | ZSET | ranking por CPU |

---

## En Grafana:

| Panel | Command | Key |
|---|---|---|
| Total RAM | GET | `system:total_ram_gb` |
| Free RAM | GET | `system:ram_free_mb` |
| RAM usada | GET | `system:ram_used_mb` |
| Gráfica RAM en tiempo | XRANGE | `container:nginx` |
| Top 5 RAM | ZRANGE | `top:ram` |
| Top 5 CPU | ZRANGE | `top:cpu` |

Para el **Top 5** en Grafana usa `ZRANGE` con `REV` activado y `Count: 5` para obtener los 5 mayores.
