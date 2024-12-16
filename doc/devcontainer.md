### Documentación de `devcontainer` para proyectos de desarrollo en VS Code

#### **¿Qué es `devcontainer` y para qué sirve?**
`devcontainer` es una característica de Visual Studio Code (VSCode) que te permite configurar un entorno de desarrollo completo y reproducible dentro de un contenedor de Docker. Esto es especialmente útil para garantizar que el entorno de desarrollo de todos los miembros de un equipo sea consistente, sin importar su sistema operativo o configuración local.

Los archivos de `devcontainer` definen cómo debe configurarse este entorno de desarrollo dentro de un contenedor, permitiendo que todas las dependencias, herramientas y configuraciones se gestionen de manera aislada.

#### **Estructura básica de `devcontainer.json`**

El archivo `devcontainer.json` define la configuración del entorno de desarrollo en contenedores. A continuación, se muestra la estructura básica de este archivo:

```json
{
  "name": "Nombre del contenedor",

  // Define el archivo de Docker Compose a usar
  "dockerComposeFile": "ruta/al/archivo/docker-compose.yml",

  // Nombre del servicio de Docker que se utilizará como contenedor de desarrollo
  "service": "nombre_del_servicio",

  // Carpeta de trabajo dentro del contenedor
  "workspaceFolder": "/ruta/de/trabajo/dentro/del/contenedor",

  // Personalización de VS Code, incluyendo extensiones y configuraciones
  "customizations": {
    "vscode": {
      "settings": {
        // Configuraciones para el lenguaje, como Go en este caso
        "go.useLanguageServer": true,
        "editor.formatOnSave": true
      },
      "extensions": [
        "golang.Go",  // Extensión de Go
        "ms-vscode-remote.remote-containers"  // Extensión de contenedores remotos
      ]
    }
  },

  // Comando a ejecutar después de crear el contenedor
  "postCreateCommand": "go mod tidy",

  // Puertos que deben estar disponibles desde el contenedor en la máquina host
  "forwardPorts": [2345],

  // Variables de entorno específicas
  "remoteEnv": {
    "GOPROXY": "https://proxy.golang.org,direct"
  },

  // Configuración de características del contenedor
  "features": {
    "ghcr.io/devcontainers/features/docker-in-docker:1": {}
  }
}
```

#### **Uso común de `devcontainer`**
1. **Reproducibilidad del entorno de desarrollo:** Garantiza que todos los desarrolladores trabajen con las mismas versiones de herramientas, dependencias y configuraciones.
2. **Integración con Docker Compose:** Puedes configurar múltiples servicios que interactúan entre sí (bases de datos, APIs, microservicios, etc.) en archivos `docker-compose.yml`.
3. **Configuraciones automáticas:** Puedes especificar extensiones, configuraciones de lenguaje y comandos de instalación que se ejecutan automáticamente cuando el contenedor se levanta.

#### **Pasos para configurar y usar un `devcontainer`**

1. **Creación del archivo `devcontainer.json`**:
   - Coloca el archivo `devcontainer.json` en el directorio `.devcontainer` dentro de tu proyecto (sin el punto inicial).
   - Configura el archivo con los servicios, extensiones y personalizaciones que desees.

2. **Levantamiento del contenedor**:
   - Abre tu proyecto en VSCode.
   - Si ya tienes Docker ejecutándose y configurado, puedes ir a la paleta de comandos (Ctrl+Shift+P) y seleccionar la opción **"Dev Containers: Reopen in Container"**.
   - El contenedor se construirá y abrirá el entorno de desarrollo dentro de él.

3. **Montaje de volúmenes**:
   - Los cambios que hagas dentro del contenedor en tu código se reflejan en tu máquina host, gracias a la opción `mounts` en el archivo `devcontainer.json`.
   
   Ejemplo:
   ```json
   "mounts": [
     "source=${localWorkspaceFolder},target=/app,type=bind,consistency=cached"
   ]
   ```

4. **Uso de Docker Compose**:
   Si tienes varios servicios definidos en un archivo `docker-compose.yml`, puedes especificar este archivo en `devcontainer.json` y definir el servicio que actuará como entorno de desarrollo:

   ```json
   "dockerComposeFile": "../config/docker-compose.dev.yml",
   "service": "nombre_del_servicio"
   ```

5. **Extensiones y configuración de VSCode**:
   Puedes preinstalar extensiones de VSCode dentro del contenedor para mejorar tu flujo de trabajo.

   ```json
   "customizations": {
     "vscode": {
       "extensions": [
         "golang.Go",
         "ms-vscode-remote.remote-containers"
       ]
     }
   }
   ```

6. **Comandos posteriores a la creación y al inicio**:
   Para asegurarte de que el entorno esté listo después de la creación del contenedor, puedes utilizar los comandos `postCreateCommand` y `postStartCommand` para ejecutar acciones como la instalación de dependencias o la ejecución de tests.

#### **Problemas comunes y soluciones**

1. **Problemas al cargar `devcontainer.json`:** Si tu archivo tiene un punto inicial (es decir, `.devcontainer.json`), VSCode no lo reconocerá. El archivo debe llamarse `devcontainer.json` sin el punto al principio.

2. **Errores de montaje de volúmenes:**
   Asegúrate de que la ruta especificada en `"dockerComposeFile"` sea correcta y accesible. También verifica que las rutas de `volumes` estén bien configuradas en tu archivo `docker-compose`.

3. **Errores de servicios no encontrados:**
   - Si aparece el error `no service selected`, asegúrate de que el nombre del servicio en `"service"` coincida exactamente con el nombre del servicio en el archivo `docker-compose.yml`.

4. **Problemas al iniciar el contenedor:**
   Si el contenedor no arranca y devuelve un error al intentar levantarlo desde `devcontainer`, verifica los logs en **View > Output** > **Dev Containers** para identificar posibles errores de configuración o permisos.

5. **Reiniciar Docker y VSCode:**
   Si tienes problemas recurrentes con el contenedor, intenta reiniciar Docker y VSCode, y borra los archivos temporales de `Dev Containers`.

#### **Resumen del problema encontrado**
Durante nuestra configuración, el problema principal fue que el archivo `devcontainer.json` tenía un punto inicial, lo que impedía que VSCode lo reconociera correctamente. Cambiar el nombre del archivo a `devcontainer.json` resolvió ese problema, pero luego surgieron otros problemas relacionados con la configuración y levantamiento del contenedor en `Dev Containers`. Decidiste no continuar con la solución de estos problemas en este momento.

#### **Conclusión**
El uso de `devcontainer` es muy potente para entornos de desarrollo consistentes y reproducibles. Sin embargo, es crucial tener una configuración adecuada de los archivos de Docker Compose, los servicios y los montajes de volúmenes para que todo funcione correctamente.