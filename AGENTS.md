# Estructura del Proyecto y Protocolo

Este archivo describe la organización visible del repositorio y las reglas de actualización.

## Estructura de Directorios (Visible)

- **docs/**: Contiene toda la base teórica.
  - **unidad_1/**: Linux, virtualización, contenedores, Go y Rust.
  - **unidad_2/**: Kubernetes y sostenibilidad.

- **ejemplos/**: Código fuente de las prácticas realizadas.
  - **semana-1/**: API básica en Go.
  - **semana-2/**: API con Fiber.

*(Nota: Las carpetas de las semanas 3-14 y otros ejemplos están actualmente ocultas por .gitignore)*

## Protocolo de Actualización de README

**Instrucción Crítica:**

Si modificas o se modifica el archivo `.gitignore` para permitir el seguimiento de una nueva carpeta (ej. liberando el contenido de una semana futura):

1. **Debes** editar inmediatamente el `README.md` en la raíz.
2. Agrega la nueva carpeta a la lista de **Ejemplos de Código** o **Documentación**.
3. Incluye una descripción breve y formal de su contenido.

El objetivo es mantener el `README.md` sincronizado con el contenido visible del repositorio en todo momento.

## Reglas Importantes

- **Prohibido usar emojis.**
