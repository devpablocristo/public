### Actualización del SDK en proyectos dependientes

Cuando tu proyecto depende de un SDK que está en constante desarrollo, y no quieres estar creando tags por cada cambio, puedes apuntar directamente a la rama `main` del SDK. Esto te asegura que siempre utilices la última versión disponible.

#### Pasos para actualizar la versión del SDK en tu proyecto:

1. **Apunta a la rama `main` en tu `go.mod`:**
   Para que tu proyecto siempre utilice la versión más reciente del SDK, en el archivo `go.mod` usa la siguiente referencia:

   ```bash
   go get github.com/devpablocristo/golang-sdk@main
   ```

2. **Actualizar a la última versión:**
   Cada vez que se haga un nuevo commit en el SDK y quieras traer los últimos cambios a tu proyecto, ejecuta el siguiente comando:

   ```bash
   go get github.com/devpablocristo/golang-sdk@main
   ```

   Esto actualizará tu archivo `go.mod` para apuntar a la última versión de la rama `main` del SDK.

#### Consideraciones:
- **Cambios automáticos**: Go no actualiza automáticamente las dependencias. Por lo tanto, siempre que haya cambios en el SDK, debes ejecutar `go get` para traer los últimos commits.
- **Control de versiones**: Esta estrategia es útil en desarrollo, pero en producción se recomienda usar versiones etiquetadas (tags) para tener más control sobre las versiones que usas.