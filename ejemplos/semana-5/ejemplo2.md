## 1. Código actualizado

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const watchDir = "/home/tu_usuario/vigilante-daemon" // <-- cambia tu_usuario por el tuyo

type fileSnapshot map[string]time.Time

func snapshot(dir string) fileSnapshot {
	snap := make(fileSnapshot)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		snap[path] = info.ModTime()
		return nil
	})
	return snap
}

func compare(old, new fileSnapshot) {
	for path, modTime := range new {
		if _, exists := old[path]; !exists {
			log.Printf("[CREADO]     %s", path)
		} else if old[path] != modTime {
			log.Printf("[MODIFICADO] %s (mod: %s)", path, modTime.Format("15:04:05"))
		}
	}
	for path := range old {
		if _, exists := new[path]; !exists {
			log.Printf("[ELIMINADO]  %s", path)
		}
	}
}

func main() {
	os.MkdirAll(watchDir, 0755)

	f, err := os.OpenFile("/var/log/vigilante-daemon.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("=== vigilante-daemon iniciado | Vigilando: %s ===", watchDir)
	fmt.Printf("Vigilando %s — revisa /var/log/vigilante-daemon.log\n", watchDir)

	prev := snapshot(watchDir)

	for {
		time.Sleep(5 * time.Second)
		curr := snapshot(watchDir)
		compare(prev, curr)
		prev = curr
	}
}
```

> **Cambia** `tu_usuario` por tu usuario real, por ejemplo `/home/juan/vigilante-daemon`  
> O si prefieres usar la variable de entorno directamente, reemplaza la constante por:
> ```go
> var watchDir = filepath.Join(os.Getenv("HOME"), "vigilante-daemon")
> ```
> y quita el `const`.

---

## 2. Crear la carpeta vigilada

```bash
mkdir ~/vigilante-daemon
```

---

## 3. Compilar con el nuevo nombre

```bash
go build -o /usr/local/bin/vigilante-daemon main.go
```

---

## 4. Archivo systemd — `/etc/systemd/system/vigilante-daemon.service`

```ini
[Unit]
Description=Daemon vigilante de ~/vigilante-daemon
After=network.target

[Service]
ExecStart=/usr/local/bin/vigilante-daemon
Restart=always
RestartSec=3
Environment=HOME=/home/tu_usuario   # <-- cambia esto también

[Install]
WantedBy=multi-user.target
```

---

## 5. Activar

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now vigilante-daemon
```

---

## 6. Probar

```bash
# Crear, modificar y borrar archivos en la carpeta vigilada
touch ~/vigilante-daemon/prueba.txt
echo "hola" >> ~/vigilante-daemon/prueba.txt
rm ~/vigilante-daemon/prueba.txt

# Ver logs en tiempo real
tail -f /var/log/vigilante-daemon.log
```

Verás:
```
=== vigilante-daemon iniciado | Vigilando: /home/juan/vigilante-daemon ===
[CREADO]     /home/juan/vigilante-daemon/prueba.txt
[MODIFICADO] /home/juan/vigilante-daemon/prueba.txt (mod: 15:10:42)
[ELIMINADO]  /home/juan/vigilante-daemon/prueba.txt
```

---

## 7. Apagar

```bash
sudo systemctl stop vigilante-daemon
sudo systemctl disable vigilante-daemon
```
