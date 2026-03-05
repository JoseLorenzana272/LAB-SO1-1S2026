Para usar el Cronjob

Ejecuta esto:

```bash
sudo apt-get update && sudo apt-get install -y cron
```

Luego verifica que esté corriendo:

```bash
sudo service cron start
sudo service cron status
```

Y confirma que `crontab` ya esté disponible:

```bash
which crontab
```


---

## Cómo cargar algo a crontab

### Formas principales:

**1. Editar interactivamente (la más común)**
```bash
crontab -e
```
Se abre un editor donde agregas tu línea y guardas.

---

**2. Agregar una línea desde terminal (sin abrir editor)**
```bash
(crontab -l 2>/dev/null; echo "* * * * * /ruta/script.sh") | crontab -
```
Esto conserva los cronjobs existentes y agrega el nuevo al final.

---

**3. Cargar desde un archivo**
```bash
crontab /ruta/mi-crontab.txt
```
⚠️ Esto **reemplaza todo** el crontab con el contenido del archivo.

---

**4. Ver los cronjobs actuales**
```bash
crontab -l
```

**5. Eliminar todos los cronjobs**
```bash
crontab -r
```

---

### Sintaxis de una línea crontab:
```
* * * * * /ruta/al/script.sh
│ │ │ │ │
│ │ │ │ └── día de la semana (0-7, donde 0 y 7 = domingo)
│ │ │ └──── mes (1-12)
│ │ └────── día del mes (1-31)
│ └──────── hora (0-23)
└────────── minuto (0-59)
```

### Ejemplos comunes:
```bash
* * * * *        # cada minuto
*/5 * * * *      # cada 5 minutos
0 * * * *        # cada hora en punto
0 2 * * *        # todos los días a las 2am
0 2 * * 1        # cada lunes a las 2am
```
