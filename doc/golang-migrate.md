# Golang-migrate: Gestión de Migraciones en Bases de Datos

**Golang-migrate** es una herramienta escrita en Go para gestionar las migraciones de bases de datos. Te permite aplicar, revertir y versionar cambios en el esquema de una base de datos de manera controlada y segura. Al utilizar `golang-migrate`, puedes asegurar que los cambios estructurales se apliquen de forma coherente en todos los entornos.

### Características Principales

1. **Compatibilidad con Múltiples Bases de Datos**:
   - Soporte para PostgreSQL, MySQL, SQLite, SQL Server, Cassandra, entre otros.
   
2. **Control de Versiones**:
   - Las migraciones se versionan, lo que permite aplicar, revertir o saltar a una versión específica.
   
3. **Integración con Múltiples Fuentes de Migración**:
   - Puedes usar migraciones desde archivos locales, AWS S3, Google Cloud Storage, entre otros.

4. **Ejecución de Migraciones**:
   - `up`: Aplica todas las migraciones pendientes.
   - `down`: Revierte las migraciones.
   - `goto`: Aplica o revierte hasta una versión específica.
   - `force`: Fuerza la base de datos a una versión específica sin aplicar ni revertir migraciones.

---

## Instalación

### Opción 1: Mediante Go

Instala `golang-migrate` directamente usando Go:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Opción 2: Binario Precompilado

Puedes descargar un binario precompilado desde la [página de releases](https://github.com/golang-migrate/migrate/releases).

Para macOS, usando Homebrew:
```bash
brew install golang-migrate
```

---

## Creación de Migraciones

Las migraciones definen los cambios estructurales en la base de datos. Cada migración consta de dos archivos:

- **`up.sql`**: Define cómo aplicar el cambio.
- **`down.sql`**: Define cómo revertir el cambio.

### Crear una Nueva Migración

Para crear un nuevo conjunto de migraciones:
```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

Este comando genera dos archivos:

- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

### Ejemplo de Archivo de Migración

**Archivo `up.sql`**: Crear la tabla `users`:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Archivo `down.sql`**: Revertir el cambio:
```sql
DROP TABLE users;
```

---

## Aplicar y Revertir Migraciones

### Aplicar Migraciones

Para aplicar todas las migraciones pendientes:
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" up
```

### Revertir Migraciones

Para revertir la última migración:
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" down 1
```

El `1` indica que solo se revierte la última migración. Puedes cambiar este número para revertir más migraciones.

### Ir a una Migración Específica

Si deseas migrar a una versión específica (hacia arriba o hacia abajo):
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" goto 3
```

### Forzar una Versión Específica

En caso de que haya inconsistencias, puedes forzar la base de datos a una versión sin ejecutar ninguna migración:
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" force 3
```

---

## Estado y Versionado de Migraciones

`golang-migrate` mantiene un registro de las migraciones aplicadas en una tabla especial (`schema_migrations`) dentro de la base de datos. Puedes verificar el estado de las migraciones con:

```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" version
```

Esto muestra la versión actual de la base de datos y las migraciones pendientes.

---

## Uso en Aplicaciones Go

Además de ejecutarlo desde la línea de comandos, también puedes integrar `golang-migrate` directamente en tu aplicación Go.

### Instalación en Go:
```bash
go get -u github.com/golang-migrate/migrate/v4
```

### Ejemplo de Código en Go

```go
package main

import (
    "log"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    m, err := migrate.New(
        "file://migrations",
        "postgres://user:password@localhost:5432/mydb?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    // Aplicar todas las migraciones pendientes
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    // Para revertir la última migración
    if err := m.Steps(-1); err != nil {
        log.Fatal(err)
    }
}
```

---

## Buenas Prácticas

1. **Versiona los Archivos de Migración junto al Código**: Mantén los archivos de migración versionados junto al código de la aplicación para asegurar que el esquema de la base de datos esté sincronizado.

2. **Aplica Migraciones en Entornos de Prueba Primero**: Antes de aplicar migraciones en producción, asegúrate de probarlas en entornos de desarrollo y pruebas.

3. **Manejo de Errores**: Siempre captura y maneja errores cuando trabajas con migraciones, especialmente en producción, para evitar problemas de inconsistencia.

4. **Deshacer Migraciones (Rollback)**: Asegúrate de tener un mecanismo de reversión en caso de que una migración falle o cause problemas inesperados en producción.

5. **Integración Continua**: Puedes integrar `golang-migrate` en tu pipeline de CI/CD para asegurar que las migraciones se apliquen automáticamente en los despliegues.

---

## Resumen

`golang-migrate` es una herramienta robusta y flexible para manejar migraciones en PostgreSQL (y otras bases de datos) de manera eficiente. Te permite gestionar de manera efectiva los cambios en el esquema de la base de datos, ya sea desde la línea de comandos o integrándolo en tu aplicación Go. Usando buenas prácticas y asegurando la consistencia en todos los entornos, puedes evitar errores y mantener tu base de datos en sincronía con el código de la aplicación.