Para crear bases de datos en PostgreSQL mediante scripts de inicialización:

Hay varias formas de solucionarlo, una de ellas es:

1. **Usar PSQL con template1**
```sql
\connect template1;
CREATE DATABASE users_db;
```

Por ejemplo:

`001_create_users_db.sql`:
```sql
\connect template1;
CREATE DATABASE users_db;
```

`002_create_tables.sql`:
```sql
\connect users_db;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    ...
);
```

`003_insert_initial_data.sql`:
```sql
\connect users_db;
INSERT INTO users (name) VALUES ('John Doe');
```

De esta manera, aseguras que:
1. Primero se crea la base de datos (conectándote a template1)
2. Los scripts posteriores se ejecutan dentro de la base de datos que acabas de crear

Si necesitas trabajar con múltiples bases de datos, puedes crear todas las que necesites en el primer script y luego conectarte a la base de datos correspondiente en cada script posterior. Por ejemplo:

`001_create_databases.sql`:
```sql
\connect template1;
CREATE DATABASE users_db;
CREATE DATABASE products_db;
```

`002_create_users_tables.sql`:
```sql
\connect users_db;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);
```

`003_create_products_tables.sql`:
```sql
\connect products_db;
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    price DECIMAL(10,2)
);
```

`004_insert_users_data.sql`:
```sql
\connect users_db;
INSERT INTO users (name) VALUES ('John Doe');
```

`005_insert_products_data.sql`:
```sql
\connect products_db;
INSERT INTO products (name, price) VALUES ('Product 1', 99.99);
```

Los scripts se ejecutarán en orden alfabético/numérico, por lo que:
1. Primero se crean todas las bases de datos
2. Luego puedes crear las tablas en cada base de datos
3. Finalmente puedes insertar los datos iniciales

Recuerda que puedes organizar tus scripts en subdirectorios dentro de `postgres-init` para mantenerlo más ordenado, por ejemplo:
```
postgres-init/
├── 001_create_databases.sql
├── users/
│   ├── 002_create_tables.sql
│   └── 004_insert_data.sql
└── products/
    ├── 003_create_tables.sql
    └── 005_insert_data.sql
```
No necesariamente deben ir en el mismo script. Puedes separarlas siempre y cuando mantengas el orden correcto de ejecución usando la numeración. Por ejemplo:

`001_create_users_db.sql`:
```sql
\connect template1;
CREATE DATABASE users_db;
```

`002_create_products_db.sql`:
```sql
\connect template1;
CREATE DATABASE products_db;
```

`003_create_users_tables.sql`:
```sql
\connect users_db;
CREATE TABLE users (...);
```

`004_create_products_tables.sql`:
```sql
\connect products_db;
CREATE TABLE products (...);
```

La ventaja de separarlas es:
- Mejor organización
- Más fácil de mantener
- Más claro qué script corresponde a cada base de datos
- Más fácil de desactivar una base de datos si temporalmente no la necesitas (simplemente renombrando o moviendo su script)

Una estructura de directorios podría ser:

```
postgres-init/
├── users/
│   ├── 001_create_users_db.sql
│   ├── 003_create_tables.sql
│   └── 005_insert_data.sql
└── products/
    ├── 002_create_products_db.sql
    ├── 004_create_tables.sql
    └── 006_insert_data.sql
```

Otra forma sera:

Es necesario agregar el prefijo alfabético a los directorios:

```
postgres-init/
├── a_users/
│   ├── 001_create_db.sql
│   ├── 002_create_tables.sql
│   └── 003_insert_data.sql
└── b_products/
    ├── 001_create_db.sql
    ├── 002_create_tables.sql
    └── 003_insert_data.sql
```

De esta manera, PostgreSQL ejecutará:
1. Todos los scripts de `a_users` en orden
2. Luego todos los scripts de `b_products` en orden

El prefijo alfabético (`a_`, `b_`) es necesario para controlar el orden de ejecución entre directorios.

O con números `1_` y `2_` en los directorios asegurarán que:

1. Primero se ejecuten todos los scripts de `1_users` en orden:
   - 001_create_db.sql
   - 002_create_tables.sql
   - 003_insert_data.sql

2. Luego todos los scripts de `2_products` en orden:
   - 001_create_db.sql
   - 002_create_tables.sql
   - 003_insert_data.sql

Esta estructura es:
- Clara y fácil de entender
- Mantiene todos los scripts de cada base de datos juntos
- El orden de ejecución es predecible
- Fácil de mantener y expandir si necesitas agregar más bases de datos (3_orders, 4_inventory, etc.)


## Documentación PostgreSQL

### Resumen:

1. Correr container
2. Entrar al container por la terminal
3. Ejecutar: `$ psql -U postgres`
4. En mis contenedores, siempre uso:
    - **Superusuario**: `admin`
    - **Contraseña**: `admin`
    - **Usuario común**: `user`
    - **Contraseña**: `user`
5. Crear usuario y darle permisos:
    ```sql
    CREATE USER admin WITH PASSWORD 'admin';
    ```
6. Convertir en superusuario:
    ```sql
    ALTER USER admin WITH SUPERUSER;
    ```


En PostgreSQL, el nombre de usuario y la contraseña por defecto pueden ser configurados al crear los contenedores de Docker. Comúnmente, se usa el usuario `postgres` como superusuario por defecto, pero se pueden agregar usuarios adicionales con privilegios específicos.

### 1. Preparación del Entorno

1. **Configura tu archivo `config/.env`**: 
   Asegúrate de que el archivo `.env` en el directorio `config` contiene las variables necesarias para configurar la base de datos PostgreSQL. Por ejemplo:

   ```env
   POSTGRES_HOST=postgres
   POSTGRES_PORT=5432
   POSTGRES_DATABASE=my_database
   POSTGRES_USERNAME=admin
   POSTGRES_PASSWORD=admin
   ```

2. **Instala PostgreSQL usando Docker Compose**: 
   Define el servicio de PostgreSQL en tu archivo `docker-compose.yml` y asegúrate de incluir las credenciales en las variables de entorno. Por ejemplo:

   ```yaml
   postgres:
     image: postgres:latest
     container_name: postgres
     environment:
       POSTGRES_USER: admin
       POSTGRES_PASSWORD: admin
       POSTGRES_DB: my_database
     ports:
       - "5432:5432"
     volumes:
       - ./postgres_data:/var/lib/postgresql/data
     networks:
       - app-network
   ```

3. **Levanta el contenedor de PostgreSQL**: 
   Ejecuta el siguiente comando para iniciar el servicio de PostgreSQL con Docker Compose:

   ```bash
   docker-compose up -d
   ```

### 2. Creación Manual de la Base de Datos y Usuarios

1. **Accede al contenedor de PostgreSQL**:
   Una vez que el contenedor esté en ejecución, puedes acceder al contenedor de PostgreSQL con el siguiente comando:

   ```bash
   docker exec -it postgres bash
   ```

2. **Accede a la consola de PostgreSQL**:
   Dentro del contenedor, usa `psql` para conectarte a PostgreSQL como el superusuario `postgres`:

   ```bash
   psql -U postgres
   ```

3. **Crear usuario y asignar permisos**:
   Puedes crear un usuario normal y un superusuario dentro de PostgreSQL con los siguientes comandos:

   - Crear un nuevo usuario:
     ```sql
     CREATE USER user WITH PASSWORD 'user';
     ```

   - Convertir a un usuario en superusuario:
     ```sql
     ALTER USER admin WITH SUPERUSER;
     ```

4. **Verificar los usuarios existentes**:
   Usa el comando `\du` para listar los usuarios actuales en PostgreSQL:

   ```sql
   \du
   ```

### 3. Flujo para Administrar Usuarios y Permisos en PostgreSQL

1. **Correr el contenedor de PostgreSQL**:
   Asegúrate de que el contenedor de PostgreSQL esté corriendo. Si no lo está, inicia el contenedor con Docker Compose.

   ```bash
   docker-compose up -d
   ```

2. **Entrar al contenedor por la terminal**:
   Usa el siguiente comando para acceder al contenedor de PostgreSQL:

   ```bash
   docker exec -it postgres bash
   ```

3. **Conectar a PostgreSQL**:
   Una vez dentro del contenedor, conéctate a PostgreSQL usando el siguiente comando:

   ```bash
   psql -U postgres
   ```

4. **Gestión de usuarios**:
   - **Superusuario**: `admin`
     ```sql
     CREATE USER admin WITH PASSWORD 'admin';
     ALTER USER admin WITH SUPERUSER;
     ```

   - **Usuario común**: `user`
     ```sql
     CREATE USER user WITH PASSWORD 'user';
     ```

5. **Verificación y gestión**:
   - **Ver usuarios**:
     ```sql
     \du
     ```

6. **Salir de PostgreSQL y del contenedor**:
   - Para salir de PostgreSQL, usa `\q`.
   - Para salir del contenedor, usa `exit`. 

Las **migraciones** en PostgreSQL son un conjunto de cambios incrementales que permiten mantener el esquema de la base de datos actualizado a lo largo del ciclo de vida de una aplicación. Esto incluye la creación y modificación de tablas, índices, columnas, restricciones, y otros objetos de la base de datos.

Una **migración** es básicamente un script que define cambios específicos en el esquema, y puede ser aplicada para actualizar la base de datos. Cada cambio tiene una dirección **"hacia adelante"** (migrar) y opcionalmente una dirección **"hacia atrás"** (rollback), que permite deshacer esos cambios si es necesario.

### ¿Cómo funcionan las migraciones en PostgreSQL?

1. **Archivo de migraciones**:
   Cada migración generalmente tiene dos partes: 
   - **Up (migrar):** la parte que define los cambios hacia adelante (p.ej., crear tablas, agregar columnas).
   - **Down (rollback):** opcionalmente, la parte que revierte esos cambios (p.ej., eliminar tablas, quitar columnas).

   Los archivos de migración suelen tener nombres que indican el orden en que deben ser aplicados, como `001_create_users_table.sql`, `002_add_email_to_users.sql`.

2. **Registro del estado de las migraciones**:
   Para llevar un control de las migraciones aplicadas, las herramientas de migración crean una tabla en la base de datos (por ejemplo, una tabla llamada `schema_migrations`) que mantiene un registro de las migraciones que ya se han aplicado, de modo que no se vuelvan a ejecutar innecesariamente.

3. **Herramientas de migración**:
   Herramientas como `golang-migrate` o `Flyway` gestionan el proceso de migración en PostgreSQL. Estas herramientas leen los archivos de migración y aplican los cambios al esquema de la base de datos. También permiten revertir migraciones en caso de errores o cambios necesarios.

   Ejemplo de herramientas populares:
   - **`golang-migrate`**: Una de las bibliotecas más populares para Go, usada para manejar migraciones en diversas bases de datos, incluyendo PostgreSQL.
   - **Flyway**: Herramienta de migración multiplataforma y multibase de datos, usada tanto en entornos de desarrollo como en producción.

4. **Ejecución de migraciones**:
   Durante la ejecución de una migración, la herramienta lee el archivo de migración, aplica las instrucciones contenidas en la parte "up" al esquema de la base de datos, y actualiza la tabla `schema_migrations` con el número o identificador de la migración aplicada.

5. **Reversión (Rollback)**:
   Si se necesita revertir una migración, se ejecuta el rollback, que aplica la parte "down" de las migraciones. Esto puede eliminar tablas o revertir cualquier otro cambio hecho anteriormente.

### Ejemplo de un archivo de migración:

#### 001_create_users_table.sql (Up)
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 001_create_users_table.sql (Down)
```sql
DROP TABLE IF EXISTS users;
```

### Flujo de una migración en Go con `golang-migrate`:

1. **Instalar `golang-migrate`**:
   Si no lo tienes instalado, puedes instalarlo con Go:
   ```bash
   go get -u -d github.com/golang-migrate/migrate/cmd/migrate
   ```

2. **Crear una nueva migración**:
   Puedes crear archivos de migración para "up" y "down":
   ```bash
   migrate create -ext sql -dir migrations -seq create_users_table
   ```

   Esto generará dos archivos:
   - `000001_create_users_table.up.sql`: Aquí defines las instrucciones de la migración.
   - `000001_create_users_table.down.sql`: Aquí defines las instrucciones para revertir la migración.

3. **Aplicar migraciones**:
   Para aplicar todas las migraciones pendientes, se ejecuta el siguiente comando:
   ```bash
   migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
   ```

4. **Revertir una migración**:
   Para deshacer la última migración, se puede ejecutar:
   ```bash
   migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
   ```

### Ventajas de usar migraciones:
- **Control de versiones** del esquema de la base de datos.
- Facilita la **sincronización** de bases de datos en diferentes entornos (desarrollo, pruebas, producción).
- Permite realizar **rollback** si una migración falla o hay que revertir cambios.

### Buenas prácticas al usar migraciones:
- **Versionar** los archivos de migración junto con el código en un sistema de control de versiones como Git.
- **Probar** las migraciones en un entorno de staging antes de aplicarlas en producción.
- **Tener rollback** para todas las migraciones críticas.

### Conclusión:
Las migraciones son fundamentales en PostgreSQL para mantener el esquema de la base de datos consistente con el código de la aplicación a lo largo del tiempo, especialmente en equipos o sistemas distribuidos. Las herramientas como `golang-migrate` facilitan la gestión de migraciones de manera automática y eficiente.   