El método de conexión que utiliza una cadena de texto para especificar todos los detalles necesarios para conectarse a un servicio se conoce comúnmente como **Connection String** o **Connection URI**. Este método se emplea para conectar aplicaciones a bases de datos, servicios de mensajería, sistemas de archivos, y más.

### Nombres Comunes:
- **Connection String**: Nombre más general para cualquier cadena de conexión.
- **Connection URI**: Cuando la cadena tiene formato de URI/URL.
- **DSN (Data Source Name)**: Especialmente en el contexto de bases de datos y ODBC.

### Componentes Comunes en una Connection String:
- **Protocolo**: Define el tipo de servicio (e.g., `mysql`, `amqp`, `postgresql`).
- **Credenciales de usuario**: Incluye usuario y contraseña (`user:password`).
- **Host**: Dirección del servidor (puede ser un nombre de dominio o IP).
- **Puerto**: Puerto donde se está ejecutando el servicio.
- **Ruta o base de datos**: Nombre del recurso o base de datos a la que se accede.

### Ejemplos:

1. **Bases de Datos Relacionales**:
   - **PostgreSQL**: 
     ```plaintext
     postgresql://user:password@localhost:5432/mydatabase
     ```
   - **MySQL**:
     ```plaintext
     mysql://user:password@localhost:3306/mydatabase
     ```

2. **Bases de Datos NoSQL**:
   - **MongoDB**:
     ```plaintext
     mongodb://user:password@localhost:27017/mydatabase
     ```
   - **Redis**:
     ```plaintext
     redis://user:password@localhost:6379/0
     ```

3. **Mensajería y Colas**:
   - **RabbitMQ**:
     ```plaintext
     amqp://user:password@localhost:5672/
     ```

4. **Sistemas de Archivos Distribuidos**:
   - **Amazon S3**:
     ```plaintext
     s3://access_key:secret_key@bucket-name.s3.amazonaws.com/
     ```

5. **LDAP (Lightweight Directory Access Protocol)**:
   - **LDAP**:
     ```plaintext
     ldap://user:password@ldap.example.com:389/
     ```

Este método de conexión es fundamental en el desarrollo de aplicaciones, ya que encapsula toda la información necesaria para que una aplicación pueda comunicarse con un servicio externo.