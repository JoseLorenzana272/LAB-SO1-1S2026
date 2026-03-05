package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func cronJobYaExiste(script string) bool {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), script)
}

// Función para crear un cronjob simple
func crearCronJob() {
	script := os.Getenv("CRON_SCRIPT")
	if script == "" {
		script = "/home/vm1/cronjob-with-daemon/create_containers.sh"
	}

	if cronJobYaExiste(script) {
		log.Printf("El cronjob con el script %s ya existe, omitiendo creacion", script)
		return
	}

	// Verificar que el script existe y tiene permisos de ejecución
	if _, err := os.Stat(script); os.IsNotExist(err) {
		log.Fatalf("El script %s no existe", script)
	}

	// Hacer el script ejecutable
	if err := os.Chmod(script, 0755); err != nil {
		log.Printf("Advertencia: no se pudieron cambiar permisos del script: %v", err)
	}

	// Agregar cronjob que se ejecuta cada minuto
	cronCommand := fmt.Sprintf("* * * * * %s", script)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("(crontab -l 2>/dev/null; echo \"%s\") | crontab -", cronCommand))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error agregando cronjob: %v, output: %s", err, string(output))
	}

	log.Printf("Cronjob creado correctamente: %s", cronCommand)

	// Verificar que se agregó correctamente
	verifyCmd := exec.Command("crontab", "-l")
	verifyOutput, err := verifyCmd.CombinedOutput()
	if err != nil {
		log.Printf("Error verificando cronjob: %v", err)
	} else {
		log.Printf("Cronjobs actuales:\n%s", string(verifyOutput))
	}
}

func detenerContenedores() string {
	cmd := exec.Command("bash", "-c", "docker stop $(docker ps -a -q)")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("ERROR: Error deteniendo contenedores: %v, output: %s", err, string(output))
	}
	return fmt.Sprintf("Contenedores detenidos: %s", string(output))
}

func eliminarContenedores() string {
	cmd := exec.Command("bash", "-c", "docker rm $(docker ps -a -q)")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("ERROR: Error eliminando contenedores: %v, output: %s", err, string(output))
	}
	return fmt.Sprintf("Contenedores eliminados: %s", string(output))
}

func main() {
	// sino se encuentra la variable de entorno, usar el valor por defecto
	logFile := os.Getenv("DAEMON_LOGFILE")
	if logFile == "" {
		logFile = "/home/vm1/cronjob-with-daemon/logs.txt"
	}

	// 1. Crear cronjob
	crearCronJob()

	// 2. Loop infinito: escribir logs y limpiar contenedores cada 2 minutos
	for {
		msgDetener := detenerContenedores()
		msgEliminar := eliminarContenedores()

		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Error abriendo archivo: %v", err)
		}

		_, _ = f.WriteString(fmt.Sprintf("%s - %s\n", time.Now().Format(time.RFC3339), msgDetener))
		_, _ = f.WriteString(fmt.Sprintf("%s - %s\n", time.Now().Format(time.RFC3339), msgEliminar))
		f.Close()

		log.Println(msgDetener)
		log.Println(msgEliminar)

		time.Sleep(2 * time.Minute) // espera 2 minutos entre iteraciones
	}
}
