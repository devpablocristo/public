# Guía de Instalación Tech House

Este documento proporciona las instrucciones para levantar el proyecto Tech House tanto en ambiente de desarrollo como en staging.

## Requisitos Previos

- Docker instalado y funcionando en tu sistema
- Make instalado en tu sistema
- Git para clonar el repositorio

## Ambiente de Desarrollo

### Construcción de Imágenes
Para construir las imágenes de Docker necesarias para el ambiente de desarrollo:

```bash
make tech-house-dev-build
```

### Levantar el Proyecto
Para iniciar todos los servicios en modo desarrollo:

```bash
make tech-house-dev-up
```

### Detener el Proyecto
Para detener todos los servicios:

```bash
make tech-house-dev-down
```

### Ver Logs
Para ver los logs de los servicios en desarrollo:

```bash
make tech-house-dev-logs
```

## Ambiente de Staging

### Construcción de Imágenes
Para construir las imágenes de Docker necesarias para el ambiente de staging:

```bash
make tech-house-stg-build
```

### Levantar el Proyecto
Para iniciar todos los servicios en modo staging:

```bash
make tech-house-stg-up
```

### Detener el Proyecto
Para detener todos los servicios:

```bash
make tech-house-stg-down
```

### Ver Logs
Para ver los logs de los servicios en staging:

```bash
make tech-house-stg-logs
```

## Estructura del Proyecto

El proyecto utiliza un Makefile con perfiles específicos para Tech House, permitiendo una fácil gestión de los diferentes ambientes (desarrollo y staging). Los comandos están organizados de la siguiente manera:

- Comandos `dev`: Utilizados para desarrollo local
- Comandos `stg`: Utilizados para el ambiente de staging

## Solución de Problemas Comunes

1. Si los servicios no inician correctamente:
   - Verifica que los puertos necesarios estén disponibles
   - Revisa los logs usando los comandos correspondientes
   - Asegúrate de que Docker esté corriendo

2. Si la construcción falla:
   - Limpia las imágenes de Docker y vuelve a intentar
   - Verifica que tengas las últimas versiones del código

## Notas Adicionales

- Los comandos de desarrollo (`dev`) están optimizados para desarrollo local y incluyen características como hot-reload
- Los comandos de staging (`stg`) están configurados para un ambiente más cercano a producción
- Todos los comandos utilizan el perfil "tech-house" específicamente

Para cualquier problema adicional, consulta la documentación del proyecto o contacta al equipo de desarrollo.

# Documentación de la API del Gestor de Clientes

Este documento proporciona información detallada sobre los endpoints disponibles en la API del Gestor de Clientes, parte del proyecto Costumer Manager.

## Configuración Base

- **Host**: localhost
- **Puerto**: 8089
- **URL Base**: `/api/v1`
- **Content-Type**: application/json

## Endpoints de Clientes

### GET /customers
Obtiene una lista de todos los clientes.

**Request**
```http
GET http://localhost:8089/api/v1/customers
```

**Respuesta (200 OK)**
```json
{
    "customers": [
        {
            "id": 1,
            "name": "string",
            "last_name": "string",
            "email": "string",
            "phone": "string",
            "age": number,
            "birth_date": "string"
        }
    ]
}
```

### GET /customers/{id}
Obtiene un cliente específico por ID.

**Request**
```http
GET http://localhost:8089/api/v1/customers/164
```

**Respuesta (200 OK)**
```json
{
    "customers": {
        "id": 164,
        "name": "string",
        "last_name": "string",
        "email": "string",
        "phone": "string",
        "age": number,
        "birth_date": "string"
    }
}
```

### POST /customers
Crea un nuevo cliente.

**Request**
```http
POST http://localhost:8089/api/v1/customers
Content-Type: application/json

{
    "name": "Homero",
    "last_name": "Simpson",
    "email": "homeros@springfield.com",
    "phone": "555123789",
    "age": 25,
    "birth_date": "1999-01-15T00:00:00Z"
}
```

**Respuesta**
- 201 Created: Cliente creado exitosamente
- 400 Bad Request: Datos de entrada inválidos

### PUT /customers/{id}
Actualiza un cliente existente.

**Request**
```http
PUT http://localhost:8089/api/v1/customers/1
Content-Type: application/json

{
    "name": "Emma",
    "last_name": "Watson",
    "email": "emma.watson@email.com",
    "phone": "8888888888",
    "age": 20,
    "birth_date": "2004-01-15T00:00:00Z"
}
```

**Respuesta**
- 200 OK: Cliente actualizado exitosamente
- 400 Bad Request: Datos de entrada inválidos
- 404 Not Found: Cliente no encontrado

### DELETE /customers/{id}
Elimina un cliente.

**Request**
```http
DELETE http://localhost:8089/api/v1/customers/176
```

**Respuesta**
- 204 No Content: Cliente eliminado exitosamente
- 404 Not Found: Cliente no encontrado

### GET /customers/kpi
Obtiene KPIs (Indicadores Clave de Rendimiento) de los clientes.

**Request**
```http
GET http://localhost:8089/api/v1/customers/kpi
```

**Respuesta (200 OK)**
```json
{
    // métricas KPI
}
```

## Endpoints de Utilidad

### GET /ping
Endpoint simple de ping para probar la conectividad de la API.

**Request**
```http
GET http://localhost:8089/api/v1/ping
```

**Respuesta (200 OK)**
```json
{
    "message": "pong"
}
```

## Estructura de Datos del Cliente

### Campos Requeridos
```json
{
    "name": "string",       // Nombre del cliente
    "last_name": "string",  // Apellido del cliente
    "email": "string",      // Correo electrónico (único)
    "phone": "string",      // Número de teléfono
    "age": number,          // Edad
    "birth_date": "string"  // Fecha de nacimiento (formato ISO 8601)
}
```

### Validaciones
- El email debe ser único en el sistema
- La fecha de nacimiento debe estar en formato ISO 8601 (ejemplo: "2004-01-15T00:00:00Z")
- La edad debe ser un número positivo
- Los campos name, last_name, email, edad y fecha de nacimiento son obligatorios

## Respuestas de Error

La API devuelve respuestas de error estandarizadas en el siguiente formato:

```json
{
    "error": {
        "code": "string",
        "message": "string",
        "details": "string"
    }
}
```

Códigos de error comunes:
- Entrada Inválida: Cuando los datos proporcionados no cumplen con el formato esperado
- Error de Validación: Cuando los datos no pasan las validaciones de negocio
- No Encontrado: Cuando el recurso solicitado no existe
- Error Interno del Servidor: Cuando ocurre un error inesperado

## Colección de Postman

Adjunta la coleccion de postman correspondiente.


# Challenge

Descripción del Desafío (con Docker y preparado para Lambda + KPI de Clientes):

El reto consiste en desarrollar una API en Golang que permita gestionar clientes, usando Docker para contenerizar la aplicación y dejando todo listo para que la API pueda ser fácilmente desplegada en AWS Lambda. Además, se debe permitir la carga de edad y fecha de nacimiento de los clientes, y se debe implementar un endpoint que devuelva KPI de los clientes, como el promedio de edad y la desviación estándar.

Requisitos del Desafío:
   - Stack Tecnológico:
   - Lenguaje: Golang.
   - Framework HTTP: Echo o Gin.
   - Base de datos: SQLite (o una base de datos en memoria para simplificar).

// INFO: Done
Contenerización: Docker.
Gestión de dependencias: Utilizar Go Modules.
Documentación de la API: Implementar Swagger para la documentación automática de la API.
Preparado para AWS Lambda: El código debe estar estructurado para facilitar su despliegue en AWS Lambda o entornos similares.

// INFO: Done
Endpoints Requeridos:
   - GET /clients: Obtiene la lista de todos los clientes desde la base de datos en SQLite o en memoria.
   - GET /clients/{id}: Obtiene los detalles de un cliente específico.
   - POST /clients: Crea un nuevo cliente. El cuerpo del request debe incluir el nombre, apellido, email, número de teléfono, edad y fecha de nacimiento.
   - PUT /clients/{id}: Actualiza un cliente existente.
   - DELETE /clients/{id}: Elimina un cliente específico.
   - GET /clients/kpi: Devuelve KPI de los clientes, tales como:
       Promedio de edad.
       Desviación estándar de la edad.

// INFO: Done
Validaciones de Entrada:
   - Los campos name, last_name, email, edad y fecha de nacimiento son obligatorios para crear un cliente.
   - Validar que el email tenga un formato correcto.
   - El número de teléfono debe ser numérico y tener un mínimo de 7 dígitos.
   - La fecha de nacimiento debe ser válida y coherente con la edad provista.
   - Manejo de errores claros y consistentes para casos como cliente no encontrado, datos inválidos, etc.

// INFO: Done
Persistencia:
Los clientes deben almacenarse localmente en una base de datos SQLite o en una estructura en memoria, permitiendo un despliegue rápido.

// INFO: Done
Contenerización con Docker:
   - Crear un Dockerfile para la aplicación Golang, con una imagen base ligera (por ejemplo, golang:alpine).
   - El contenedor debe exponer el puerto en el que la API estará escuchando.
   - Se debe poder ejecutar la aplicación localmente usando Docker para pruebas y desarrollo.

// TODO
Preparación para AWS Lambda:
 - El código debe estar preparado para ser empaquetado y desplegado en un entorno como AWS Lambda. Se debe incluir:
    - Handler adaptado para Lambda (lambda.HandlerFunc) o estructura similar.
    - Instrucciones claras para empaquetar la aplicación en un archivo ZIP compatible con AWS Lambda.

// INFO: Done
KPI de Clientes:
Crear un endpoint GET /clients/kpi que calcule y devuelva los siguientes indicadores clave:
      - Promedio de edad: El cálculo del promedio de edad entre todos los clientes registrados.
  - Desviación estándar de edad: Calcular la variación de las edades respecto al promedio.
  - Los KPI deben calcularse en tiempo real basados en los clientes almacenados.

// INFO: Done
Testing:
   - Escribir pruebas unitarias para los handlers (endpoints).
   - Crear pruebas de integración para asegurarse de que la API funciona correctamente.

// TODO
Estructura del Proyecto:
/handlers: Los handlers o controladores para manejar las peticiones HTTP.
/models: Las estructuras de datos y lógica de negocio para los clientes.
/middleware: Cualquier middleware necesario, como autenticación.
/database: Configuración de la base de datos SQLite o lógica de persistencia en memoria.
/docker: Archivos relacionados con Docker, como el Dockerfile.
/lambda: Lógica o configuraciones necesarias para empaquetar y ejecutar la aplicación en AWS Lambda.

// NOTE: Revisar
Requisitos Técnicos Específicos:
  - Escalabilidad: El código debe estar modularizado y preparado para escalar, tanto a nivel local como en entornos serverless.
  - Optimización: Se valorará la optimización para un entorno serverless, como minimizar el tamaño del contenedor y reducir tiempos de arranque.

// TODO
Documentación:
   - Implementar Swagger para la documentación de los endpoints.
   - Incluir documentación sobre cómo ejecutar la aplicación localmente con Docker, y cómo empaquetarla para su despliegue en AWS Lambda.


Duración Estimada:
4-6 horas para un desarrollador con experiencia en Golang y Docker.


Entregables:
Código fuente del proyecto.
Dockerfile.
Instrucciones claras sobre cómo levantar el proyecto localmente con Docker y cómo empaquetarlo para AWS Lambda.
Documentación de la API (preferiblemente en Swagger).
