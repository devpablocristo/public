### Documentación: Uso de ngrok con Docker Compose y Configuración en Go

#### Opción 1: Usar ngrok solo con `docker-compose.yml` y `ngrok.yml`

Esta opción es la más sencilla y no requiere cambios en el código de tu aplicación Go. Aquí, la configuración de ngrok se gestiona completamente a través de Docker Compose y un archivo de configuración `ngrok.yml`.

**Paso 1: Crear el archivo `Dockerfile`**

```dockerfile
# Utiliza una imagen base oficial de Golang para construir la aplicación
FROM golang:1.22.3 as builder

WORKDIR /app

# Copia todos los archivos al contenedor
COPY . .

# Construye la aplicación
RUN go build -o main .

# Utiliza una imagen base más ligera para ejecutar la aplicación
FROM alpine:latest

WORKDIR /app

# Copia el binario construido desde la fase anterior
COPY --from=builder /app/main .

# Asegúrate de que el binario tenga permisos de ejecución
RUN chmod +x ./main

# Expone el puerto 8080
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
```

**Paso 2: Crear el archivo `docker-compose.yml`**

```yaml
version: "3.8"

services:
  rest:
    container_name: nimcin7
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    networks:
      - app-network
    restart: on-failure

  ngrok:
    image: ngrok/ngrok:latest
    container_name: ngrok
    command: ["start", "--all", "--config", "/etc/ngrok.yml"]
    volumes:
      - ./ngrok.yml:/etc/ngrok.yml
    ports:
      - "5050:5050"
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge
```

**Paso 3: Crear el archivo `ngrok.yml`**

```yaml
version: "2"
authtoken: YOUR_NGROK_AUTHTOKEN
web_addr: "127.0.0.1:5050"
tunnels:
  my-tunnel:
    proto: http
    addr: rest:8080
    domain: brave-dane-forcibly.ngrok-free.app
```

**Paso 4: Crear el archivo `main.go`**

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Gin and ngrok!")
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting Gin server: %v", err)
	}
}
```

**Paso 5: Construir y Levantar los Servicios**

Navega al directorio del proyecto y ejecuta:

```sh
docker-compose up --build
```

**Paso 6: Verificar los Logs**

Verifica los logs del contenedor `ngrok` para asegurarte de que se está utilizando el dominio especificado:

```sh
docker logs ngrok
```

Con esta configuración, ngrok estará gestionado completamente por Docker Compose y `ngrok.yml`, y tu aplicación Go no necesitará ningún cambio adicional para usar ngrok.

#### Opción 2: Controlar ngrok desde el Código Go

Si deseas un mayor control sobre ngrok y prefieres gestionarlo desde el código de tu aplicación Go, puedes usar la función `ngrok.Listen` para iniciar y detener túneles programáticamente. 

**Ejemplo Completo con Integración Directa en Go**

**Paso 1: Crear el archivo `Dockerfile`**

Este archivo sigue siendo el mismo que en la Opción 1.

**Paso 2: Crear el archivo `docker-compose.yml`**

Este archivo también sigue siendo el mismo que en la Opción 1, excepto que no necesitas el servicio `ngrok` si controlas ngrok desde el código Go.

**Paso 3: Crear el archivo `main.go`**

```go
package main

import (
	"context"
	"log"
	"github.com/gin-gonic/gin"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Gin and ngrok!")
	})

	go func() {
		if err := runNgrok(context.Background()); err != nil {
			log.Fatalf("Error starting ngrok: %v", err)
		}
	}()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting Gin server: %v", err)
	}
}

func runNgrok(ctx context.Context) error {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(config.WithForwardsTo(":8080")),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	log.Println("Ingress established at:", listener.URL())

	return nil
}
```

**Paso 4: Construir y Levantar los Servicios**

Navega al directorio del proyecto y ejecuta:

```sh
docker-compose up --build
```

Con esta configuración, ngrok es controlado directamente desde tu código Go, lo que te permite mayor flexibilidad y control sobre cómo y cuándo se establecen los túneles ngrok.

### Conclusión

Ambas opciones tienen sus ventajas. Usar Docker Compose y `ngrok.yml` es simple y efectivo para la mayoría de los casos. Sin embargo, si necesitas un control más fino y dinámico, integrar ngrok directamente en tu código Go te da esa flexibilidad. Elige la opción que mejor se adapte a tus necesidades y flujo de trabajo.

La función `runNgrok` en Go es útil cuando deseas integrar y controlar la configuración y el ciclo de vida de ngrok directamente desde tu aplicación Go, en lugar de depender únicamente de una configuración externa a través de Docker Compose y el archivo `ngrok.yml`. Aquí hay algunas razones por las que podrías optar por usar la integración directa en Go en lugar de o además de la configuración en Docker Compose:

1. **Control Programático**: Puedes iniciar y detener túneles ngrok programáticamente desde tu aplicación, lo que te permite tener un control más fino sobre cuándo y cómo se establecen los túneles.
2. **Configuración Dinámica**: Puedes cambiar la configuración de ngrok en tiempo de ejecución basado en ciertas condiciones o configuraciones de la aplicación.
3. **Simplificación de la Configuración**: Para entornos donde no puedes o no quieres usar Docker Compose, integrar ngrok directamente en tu aplicación puede ser una solución más sencilla.
4. **Desarrollo Local**: Durante el desarrollo local, puedes iniciar ngrok directamente desde tu aplicación sin necesidad de ejecutar Docker Compose.

Si decides usar Docker Compose y `ngrok.yml` para configurar ngrok, no necesitas la función `runNgrok` en tu aplicación Go. Sin embargo, si prefieres tener un control más directo desde el código de tu aplicación, puedes optar por usar esa función.












La doc pertenece a otra API pero se uso esa misma api para crear la config de esta, igual revisar y ajustar.

### Configuración de Ngrok, Docker, Docker-Compose y Golang, utilizando un dominio específico

#### Estructura del Proyecto

```
/project-root
  - Dockerfile
  - docker-compose.yml
  - ngrok.yml
  - main.go
```

### Archivo Dockerfile

```dockerfile
# Utiliza una imagen base oficial de Golang para construir la aplicación
FROM golang:1.22.3 as builder

WORKDIR /app

# Copia todos los archivos al contenedor
COPY . .

# Construye la aplicación
RUN go build -o main .

# Utiliza una imagen base más ligera para ejecutar la aplicación
FROM alpine:latest

WORKDIR /app

# Copia el binario construido desde la fase anterior
COPY --from=builder /app/main .

# Asegúrate de que el binario tenga permisos de ejecución
RUN chmod +x ./main

# Expone el puerto 8080
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
```

### Archivo docker-compose.yml

```yaml
version: "3.8"

services:
  rest:
    container_name: nimcin7
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    networks:
      - app-network
    restart: on-failure

  ngrok:
    image: ngrok/ngrok:latest
    container_name: ngrok
    command: ["start", "--all", "--config", "/etc/ngrok.yml"]
    volumes:
      - ./ngrok.yml:/etc/ngrok.yml
    ports:
      - "5050:5050"
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge
```

### Archivo ngrok.yml

```yaml
version: "2"
authtoken: YOUR_NGROK_AUTHTOKEN
web_addr: "127.0.0.1:5050"
tunnels:
  my-tunnel:
    proto: http
    addr: rest:8080
    domain: brave-dane-forcibly.ngrok-free.app
```

### Archivo main.go

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {
	// Iniciar el servidor HTTP de forma concurrente
	go func() {
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Iniciar ngrok para exponer el puerto 8080
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(config.WithForwardsTo(":8080")),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	log.Println("Ingress established at:", listener.URL())

	return http.Serve(listener, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from ngrok-go!")
}
```

### Pasos para Ejecutar

1. **Construir y Levantar los Servicios**:
   Navega al directorio del proyecto y ejecuta:

   ```sh
   docker-compose up --build
   ```

2. **Verificar los Logs**:
   Verifica los logs del contenedor `ngrok` para asegurarte de que se está utilizando el dominio especificado:

   ```sh
   docker logs ngrok
   ```

### Detalles Adicionales

- Asegúrate de reemplazar `YOUR_NGROK_AUTHTOKEN` en `ngrok.yml` con tu token de autenticación de Ngrok.
- La configuración de `web_addr` en `ngrok.yml` permite que la interfaz web de Ngrok esté disponible en el puerto 5050.
- Puedes acceder a la interfaz web de Ngrok para inspeccionar el tráfico HTTP en `http://localhost:5050`.
- El dominio especificado (`brave-dane-forcibly.ngrok-free.app`) se utilizará para exponer tu servicio.
- La URL pública también puede obtenerse iniciando sesión en tu cuenta de Ngrok y revisando los túneles activos.

### Verificación

Para asegurarte de que todo está funcionando correctamente, sigue estos pasos:

1. Abre un navegador web y navega a la URL proporcionada por ngrok. Esta URL se mostrará en los logs del contenedor `ngrok`.
2. Asegúrate de que la aplicación responda correctamente mostrando "Hello from ngrok-go!".

### Conclusión

Esta configuración te permite exponer tu aplicación Golang a través de Ngrok utilizando un dominio específico, todo orquestado con Docker y Docker-Compose. Esto facilita el desarrollo y las pruebas en un entorno local con acceso a una URL pública personalizada.