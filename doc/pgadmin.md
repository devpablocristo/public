## Documentación pgAdmin

### Conexión con PostgreSQL

pgAdmin no se conecta automáticamente a tu servidor PostgreSQL. Necesitas configurarlo manualmente. Aquí te explico cómo hacerlo:

#### Paso 1: Acceder a pgAdmin

1. Abre pgAdmin en tu navegador en `http://localhost:8081` (o el puerto que hayas configurado).

#### Paso 2: Iniciar Sesión en pgAdmin

1. Usa las credenciales configuradas en el archivo `docker-compose.yml`:
   - **Correo electrónico**: admin@admin.com
   - **Contraseña**: admin

#### Paso 3: Configurar una Nueva Conexión al Servidor PostgreSQL

1. **Añadir un nuevo servidor**:
   - En el panel de la izquierda, haz clic derecho en "Servers" y selecciona "Create" -> "Server...".

2. **Configurar la conexión del servidor**:
   - En la pestaña "General":
     - **Name**: `Postgres Local`
   - En la pestaña "Connection":
     - **Host name/address**: `postgres` (nombre del servicio en `docker-compose.yml`).
     - **Port**: `5432`.
     - **Maintenance database**: `my_db`.
     - **Username**: `admin`.
     - **Password**: `admin`.

3. **Guardar la configuración**:
   - Haz clic en "Save".

#### Paso 4: Verificar la Conexión

1. Una vez guardada la configuración, deberías ver tu nuevo servidor en el panel de la izquierda.
2. Haz clic en él para expandir y ver las bases de datos y otros objetos dentro de PostgreSQL.

### Ejemplo Visual de Configuración

#### General

- **Name**: `Postgres Local`

#### Connection

- **Host name/address**: `postgres`
- **Port**: `5432`
- **Maintenance database**: `my_db`
- **Username**: `admin`
- **Password**: `admin`

### Instrucciones para Crear la Base de Datos y la Tabla en pgAdmin

#### Paso 1: Crear la Base de Datos

1. En pgAdmin, abre el Query Tool y ejecuta:
   ```sql
   CREATE DATABASE my_db;
   ```

#### Paso 2: Crear el Usuario y Otorgar Permisos

1. En el Query Tool, ejecuta:
   ```sql
   CREATE USER admin WITH PASSWORD 'admin';
   GRANT ALL PRIVILEGES ON DATABASE my_db TO admin;
   ```

#### Paso 3: Crear la Tabla

1. Conéctate a `my_db` en el Query Tool y ejecuta:
   ```sql
   CREATE TABLE events (
       id UUID PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       description TEXT,
       location VARCHAR(255),
       start_time TIMESTAMP NOT NULL,
       end_time TIMESTAMP,
       category VARCHAR(50),
       creator_id UUID NOT NULL,
       is_public BOOLEAN NOT NULL DEFAULT true,
       is_recurring BOOLEAN NOT NULL DEFAULT false,
       series_id UUID,
       status VARCHAR(50) NOT NULL
   );
   ```

### Notas Importantes

- Sigue estos pasos en orden para evitar errores. Crear la base de datos, el usuario y la tabla en una sola ejecución puede causar errores debido a restricciones de transacciones en PostgreSQL.

## Ver las Tablas de una Base de Datos en pgAdmin

### Pasos

1. Abre pgAdmin y conéctate al servidor.
2. Expande la base de datos `my_db`.
3. Expande `Schemas` -> `public` -> `Tables` para ver todas las tablas.

Siguiendo estos pasos, podrás ver y gestionar las tablas en tu base de datos usando pgAdmin.


### Logging

pgAdmin muestra logs de forma excesiva. Para silenciarlos, realicé los siguientes pasos:

Para silenciar los logs de pgAdmin (o de cualquier servicio en un contenedor Docker), puedes configurar el nivel de logging o redirigir los logs a `/dev/null`. Existen varias opciones para hacerlo, utilice la siguiente:

#### Redirigir los logs en Docker

Puedes redirigir la salida del contenedor a `/dev/null`. Esto se puede hacer modificando el archivo `docker-compose.yml` o los comandos de `docker run`.

**Usando `docker-compose.yml`**

Edita tu archivo `docker-compose.yml` para que redirija los logs a `/dev/null`:

```yaml
version: '3.7'
services:
  pgadmin:
    image: dpage/pgadmin4
    environment:
      PGADMIN_DEFAULT_EMAIL: "admin@example.com"
      PGADMIN_DEFAULT_PASSWORD: "admin"
    ports:
      - "8081:80"
    logging:
      driver: "none" # <--- esta configuración
```
### Otras opciomes

1. Configurar los niveles de logging si es posible desde pgAdmin.
2. Configurar Nginx para desactivar los logs de acceso.
3. Configurar el logging globalmente en Docker.

### Configurar PgAdmin desde afuera de PdAdmin

Par configurar **pgAdmin** utilizando **Docker Compose** de manera que, al iniciar los contenedores, **pgAdmin** ya esté automáticamente conectado a tu servidor de **PostgreSQL** sin necesidad de configurarlo manualmente cada vez. A continuación, te proporcionaré una guía detallada para lograr esto.

## **1. Estructura del Proyecto**

Organiza tu proyecto con la siguiente estructura de directorios:

```
my_project/
├── docker-compose.yml
├── pgadmin/
│   └── servers.json
└── data/
    └── postgres/    # Volumen para datos persistentes de PostgreSQL
```

- **`docker-compose.yml`**: Define los servicios de PostgreSQL y pgAdmin.
- **`pgadmin/servers.json`**: Configuración para que pgAdmin se conecte automáticamente al servidor de PostgreSQL.
- **`data/postgres/`**: Volumen persistente para almacenar los datos de PostgreSQL.

## **2. Crear el Archivo `docker-compose.yml`**

Este archivo define los servicios de **PostgreSQL** y **pgAdmin**, junto con sus configuraciones respectivas.

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres
    restart: always
    environment:
      POSTGRES_USER: your_user
      POSTGRES_PASSWORD: your_password
      POSTGRES_DB: users_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  pgadmin:
    image: dpage/pgadmin4:6.21
    container_name: pgadmin
    restart: always
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: admin_password
      PGADMIN_LISTEN_PORT: 80
      PGADMIN_CONFIG_SERVER_MODE: "False"  # Ejecutar en modo desktop para cargar servidores.json
    ports:
      - "8080:80"
    volumes:
      - pgadmin_data:/var/lib/pgadmin
      - ./pgadmin/servers.json:/pgadmin4/servers.json  # Montar el archivo de configuración
    depends_on:
      - postgres

volumes:
  postgres_data:
  pgadmin_data:
```

### **Explicación de la Configuración**

- **Servicio `postgres`:**
  - **`image`**: Utiliza la imagen oficial de PostgreSQL versión 15.
  - **`container_name`**: Nombre del contenedor para fácil referencia.
  - **`environment`**: Configuración de la base de datos (usuario, contraseña y nombre de la base de datos).
  - **`ports`**: Mapea el puerto 5432 del contenedor al puerto 5432 de tu máquina local.
  - **`volumes`**: Monta un volumen persistente para almacenar los datos de PostgreSQL.

- **Servicio `pgadmin`:**
  - **`image`**: Utiliza la imagen oficial de pgAdmin 4 versión 6.21.
  - **`container_name`**: Nombre del contenedor para fácil referencia.
  - **`environment`**:
    - **`PGADMIN_DEFAULT_EMAIL`**: Correo electrónico para acceder a pgAdmin.
    - **`PGADMIN_DEFAULT_PASSWORD`**: Contraseña para acceder a pgAdmin.
    - **`PGADMIN_LISTEN_PORT`**: Puerto en el que pgAdmin escuchará dentro del contenedor (por defecto es 80).
    - **`PGADMIN_CONFIG_SERVER_MODE: "False"`**: Ejecuta pgAdmin en modo desktop, permitiendo cargar la configuración de servidores desde un archivo.
  - **`ports`**: Mapea el puerto 80 del contenedor al puerto 8080 de tu máquina local. Accede a pgAdmin en `http://localhost:8080`.
  - **`volumes`**:
    - **`pgadmin_data`**: Almacena la configuración y datos de pgAdmin.
    - **`./pgadmin/servers.json:/pgadmin4/servers.json`**: Monta tu archivo de configuración `servers.json` en el contenedor.
  - **`depends_on`**: Asegura que el servicio de PostgreSQL se inicie antes que pgAdmin.

## **3. Crear el Archivo `servers.json`**

Este archivo contiene la configuración del servidor de PostgreSQL que pgAdmin debe conectar automáticamente al iniciar.

### **a. Crear el Directorio para pgAdmin**

Dentro de tu directorio de proyecto, crea una carpeta llamada `pgadmin`:

```bash
mkdir pgadmin
```

### **b. Crear y Configurar `servers.json`**

Crea un archivo llamado `servers.json` dentro de la carpeta `pgadmin` con el siguiente contenido:

```json
[
  {
    "Name": "PostgreSQL Server",
    "Group": "Servers",
    "Host": "postgres",
    "Port": 5432,
    "MaintenanceDB": "users_db",
    "Username": "your_user",
    "Password": "your_password",
    "SSLMode": "prefer"
  }
]
```

### **Explicación de los Campos:**

- **`Name`**: Nombre que aparecerá en pgAdmin para el servidor.
- **`Group`**: Grupo al que pertenece el servidor en pgAdmin.
- **`Host`**: Nombre del servicio de PostgreSQL definido en `docker-compose.yml` (`postgres`).
- **`Port`**: Puerto en el que PostgreSQL está escuchando (5432).
- **`MaintenanceDB`**: Base de datos de mantenimiento (`users_db`).
- **`Username`**: Usuario de PostgreSQL (`your_user`).
- **`Password`**: Contraseña de PostgreSQL (`your_password`).
- **`SSLMode`**: Modo SSL para la conexión (`prefer`, `require`, etc.).

**Nota Importante:** Asegúrate de que las credenciales en `servers.json` coincidan con las definidas en el archivo `docker-compose.yml`.

## **4. Configurar Variables de Entorno de Forma Segura (Opcional pero Recomendado)**

Para evitar exponer directamente las credenciales en `docker-compose.yml` y `servers.json`, puedes utilizar un archivo `.env`.

### **a. Crear un Archivo `.env`**

En el directorio raíz de tu proyecto, crea un archivo llamado `.env`:

```bash
touch .env
```

### **b. Definir las Variables de Entorno en `.env`**

Añade las siguientes líneas al archivo `.env`:

```env
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=users_db
PGADMIN_DEFAULT_EMAIL=admin@example.com
PGADMIN_DEFAULT_PASSWORD=admin_password
PGADMIN_SERVERS_JSON_PASSWORD=your_password
```

### **c. Actualizar `docker-compose.yml` para Usar Variables de Entorno**

Modifica tu archivo `docker-compose.yml` para usar las variables definidas en `.env`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres
    restart: always
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  pgadmin:
    image: dpage/pgadmin4:6.21
    container_name: pgadmin
    restart: always
    environment:
      PGADMIN_DEFAULT_EMAIL: ${PGADMIN_DEFAULT_EMAIL}
      PGADMIN_DEFAULT_PASSWORD: ${PGADMIN_DEFAULT_PASSWORD}
      PGADMIN_LISTEN_PORT: 80
      PGADMIN_CONFIG_SERVER_MODE: "False"
    ports:
      - "8080:80"
    volumes:
      - pgadmin_data:/var/lib/pgadmin
      - ./pgadmin/servers.json:/pgadmin4/servers.json
    depends_on:
      - postgres

volumes:
  postgres_data:
  pgadmin_data:
```

### **d. Actualizar `servers.json` para Usar Variables de Entorno**

Desafortunadamente, **pgAdmin** no soporta directamente la interpolación de variables de entorno en archivos JSON. Sin embargo, puedes utilizar herramientas como `envsubst` para reemplazar las variables antes de iniciar los contenedores. Alternativamente, puedes mantener las credenciales en el archivo `servers.json` pero asegurarte de que el archivo `.env` esté protegido (añadiéndolo a `.gitignore`).

**Ejemplo de `servers.json` con Credenciales:**

```json
[
  {
    "Name": "PostgreSQL Server",
    "Group": "Servers",
    "Host": "postgres",
    "Port": 5432,
    "MaintenanceDB": "users_db",
    "Username": "your_user",
    "Password": "your_password",
    "SSLMode": "prefer"
  }
]
```

**Recomendación:** Mantén el archivo `.env` y `servers.json` fuera de tu control de versiones añadiendo las siguientes líneas a tu archivo `.gitignore`:

```gitignore
.env
pgadmin/servers.json
```

## **5. Iniciar los Servicios con Docker Compose**

Desde el directorio raíz de tu proyecto, ejecuta el siguiente comando para iniciar los servicios:

```bash
docker-compose up -d
```

### **Explicación del Comando:**

- **`docker-compose up`**: Crea e inicia los contenedores definidos en `docker-compose.yml`.
- **`-d`**: Ejecuta los contenedores en segundo plano (modo "detached").

## **6. Verificar la Conexión Automática en pgAdmin**

1. **Abrir pgAdmin en el Navegador:**

   Navega a [http://localhost:8080](http://localhost:8080) en tu navegador web.

2. **Iniciar Sesión:**

   Utiliza las credenciales definidas en `.env` o directamente en `docker-compose.yml`:

   - **Email**: `admin@example.com`
   - **Contraseña**: `admin_password`

3. **Verificar la Conexión al Servidor de PostgreSQL:**

   Deberías ver automáticamente el servidor de PostgreSQL ("PostgreSQL Server") en la interfaz de pgAdmin bajo el grupo "Servers".

   ![pgAdmin Server](https://www.pgadmin.org/static/img/screenshots/pgadmin4_dashboard.png)

   *Nota: La apariencia puede variar según la versión de pgAdmin.*

## **7. Consideraciones Adicionales**

### **a. Persistencia de Datos**

- **PostgreSQL**: Los datos de la base de datos se almacenan en el volumen `postgres_data`. Esto asegura que tus datos no se pierdan al reiniciar o eliminar los contenedores.
- **pgAdmin**: La configuración y los datos de pgAdmin se almacenan en el volumen `pgadmin_data`.

### **b. Seguridad**

- **Credenciales Seguras**: Asegúrate de utilizar contraseñas seguras para `POSTGRES_PASSWORD` y `PGADMIN_DEFAULT_PASSWORD`.
- **Acceso Limitado**: Considera limitar el acceso a pgAdmin únicamente a redes confiables o implementa un proxy inverso con autenticación adicional y HTTPS.

### **c. Uso de Variables de Entorno para `servers.json` (Opcional)**

Si deseas automatizar la inyección de variables de entorno en `servers.json`, puedes utilizar scripts de inicio o herramientas como `envsubst`. Sin embargo, esto añade complejidad adicional y requiere configurar scripts personalizados.

### **d. Actualización de Servicios**

Para actualizar los servicios (por ejemplo, después de modificar `servers.json`), ejecuta:

```bash
docker-compose down
docker-compose up -d
```

### **e. Solución de Problemas**

- **Logs de Contenedores**: Si encuentras problemas, revisa los logs de los contenedores para identificar errores.

  ```bash
  docker-compose logs pgadmin
  docker-compose logs postgres
  ```

- **Verificar Conexión**: Asegúrate de que el servicio de PostgreSQL esté funcionando y escuchando en el puerto correcto.

## **8. Ejemplo Completo del Archivo `docker-compose.yml`**

Aquí tienes un resumen completo del archivo `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres
    restart: always
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  pgadmin:
    image: dpage/pgadmin4:6.21
    container_name: pgadmin
    restart: always
    environment:
      PGADMIN_DEFAULT_EMAIL: ${PGADMIN_DEFAULT_EMAIL}
      PGADMIN_DEFAULT_PASSWORD: ${PGADMIN_DEFAULT_PASSWORD}
      PGADMIN_LISTEN_PORT: 80
      PGADMIN_CONFIG_SERVER_MODE: "False"
    ports:
      - "8080:80"
    volumes:
      - pgadmin_data:/var/lib/pgadmin
      - ./pgadmin/servers.json:/pgadmin4/servers.json
    depends_on:
      - postgres

volumes:
  postgres_data:
  pgadmin_data:
```

## **9. Resumen**

Siguiendo los pasos anteriores, podrás configurar **pgAdmin** para que se conecte automáticamente a tu servidor de **PostgreSQL** al iniciar los contenedores con **Docker Compose**. Esto facilita la administración de tu base de datos sin necesidad de configuraciones manuales cada vez que inicies pgAdmin.

### **Pasos Clave:**

1. **Configurar `docker-compose.yml`**: Define los servicios de PostgreSQL y pgAdmin con sus configuraciones respectivas.
2. **Crear `servers.json`**: Define la conexión al servidor de PostgreSQL para que pgAdmin lo cargue automáticamente.
3. **Montar `servers.json` en pgAdmin**: Asegura que pgAdmin cargue la configuración al iniciar.
4. **Iniciar los Servicios**: Utiliza `docker-compose up -d` para levantar los contenedores.
5. **Verificar la Conexión**: Accede a pgAdmin en `http://localhost:8080` y confirma que el servidor de PostgreSQL está conectado.

### **Consejos Finales:**

- **Mantén tus credenciales seguras**: Utiliza archivos `.env` y evita exponer contraseñas en archivos de configuración versionados.
- **Monitorea los Logs**: Siempre revisa los logs de los contenedores si encuentras problemas.
- **Actualiza Regularmente**: Mantén las imágenes de Docker actualizadas para aprovechar las últimas mejoras y parches de seguridad.

Si tienes alguna otra pregunta o necesitas más ayuda, ¡no dudes en preguntar!