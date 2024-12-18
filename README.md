# Customer Manager API

API en Golang para la gestión de clientes, desarrollada como parte del desafío Tech House.

## Índice
1. [Stack Tecnológico](#stack-tecnológico)
2. [Guía de Instalación](#guía-de-instalación)
   - [Requisitos Previos](#requisitos-previos)
   - [Ambiente de Desarrollo](#ambiente-de-desarrollo)
   - [Ambiente de Staging](#ambiente-de-staging)
3. [Documentación de la API](#documentación-de-la-api)
   - [Configuración Base](#configuración-base)
   - [Endpoints](#endpoints)
4. [Especificaciones Técnicas](#especificaciones-técnicas)
   - [Estructura del Proyecto](#estructura-del-proyecto)
   - [Modelos de Datos](#modelos-de-datos)
   - [Validaciones](#validaciones)
5. [Testing y Desarrollo](#testing-y-desarrollo)
   - [Ejecutar Pruebas](#ejecutar-pruebas) 
   - [Postman Collection](#postman-collection)
   - [Solución de Problemas Comunes](#solución-de-problemas-comunes)
6. [Despliegue en AWS Lambda](#despliegue-en-aws-lambda)
7. [Documentación con Swagger](#documentación-con-swagger)
   - [Generación de archivos de Swagger](#generación-de-archivos-de-swagger)
8. [Challenge](#challenge)

## Stack Tecnológico
- **Lenguaje**: Golang
- **Framework HTTP**: Gin  
- **Base de datos**: SQLite
- **Gestión de dependencias**: Go Modules
- **Contenerización**: Docker

## Guía de Instalación

### Requisitos Previos
- Docker instalado y funcionando
- Make instalado 
- Git para clonar el repositorio

### Ambiente de Desarrollo

#### Construcción de Imágenes
```bash
make tech-house-dev-build
```

#### Levantar el Proyecto
```bash 
make tech-house-dev-up  
```

#### Detener el Proyecto
```bash
make tech-house-dev-down
```

#### Ver Logs
```bash 
make tech-house-dev-logs
```

### Ambiente de Staging

#### Construcción de Imágenes
```bash
make tech-house-stg-build  
```

#### Levantar el Proyecto 
```bash
make tech-house-stg-up
```

#### Detener el Proyecto
```bash
make tech-house-stg-down
```

#### Ver Logs
```bash
make tech-house-stg-logs  
```

## Documentación de la API

### Configuración Base
- **Host**: localhost 
- **Puerto**: 8089
- **URL Base**: `/api/v1`
- **Content-Type**: application/json

### Endpoints

#### GET /customers
Obtiene lista de todos los clientes.

**Request**
```http
GET http://localhost:8089/api/v1/customers
```

**Response (200 OK)**   
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

#### GET /customers/{id}
Obtiene un cliente específico.

**Request**
```http
GET http://localhost:8089/api/v1/customers/164
```

**Response (200 OK)**
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

#### POST /customers
Crea un nuevo cliente.  

**Request**
```http
POST http://localhost:8089/api/v1/customers
Content-Type: application/json

{
    "name": "Homero",
    "last_name": "Simpson",
    "email": "homero@springfield.com",
    "phone": "555123789", 
    "age": 25,
    "birth_date": "1999-01-15T00:00:00Z"
}
```

**Response**
- 201 Created: Cliente creado exitosamente
- 400 Bad Request: Datos de entrada inválidos

#### PUT /customers/{id}  
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

**Response**
- 200 OK: Cliente actualizado exitosamente  
- 400 Bad Request: Datos de entrada inválidos
- 404 Not Found: Cliente no encontrado

#### DELETE /customers/{id}
Elimina un cliente.

**Request**
```http 
DELETE http://localhost:8089/api/v1/customers/176
```

**Response** 
- 204 No Content: Cliente eliminado exitosamente
- 404 Not Found: Cliente no encontrado

#### GET /customers/kpi
Obtiene métricas KPI de clientes.

**Request**
```http
GET http://localhost:8089/api/v1/customers/kpi
```

**Response (200 OK)**
```json
{
    "average_age": 35.5,
    "age_std_deviation": 7.8
}
```

## Especificaciones Técnicas

### Estructura del Proyecto
```
/
├── internal/
│   ├── config/         # Configuraciones
│   └── customer/       # Módulo de clientes  
│       ├── adapters/   # Adaptadores (inbound/outbound)
│       ├── core/       # Lógica de negocio
│       └── ports/      # Interfaces
├── pkg/                # Paquetes compartidos
└── docker/            # Archivos Docker
```

### Modelos de Datos

#### Cliente
```json
{
    "name": "string",       // Mínimo 2 caracteres, máximo 100 
    "last_name": "string",  // Mínimo 2 caracteres, máximo 100
    "email": "string",      // Formato email válido, único
    "phone": "string",      // Mínimo 7 caracteres  
    "age": number,          // Entre 1 y 150
    "birth_date": "string"  // Formato ISO 8601
}
```

### Validaciones

- **Nombre y Apellido**:  
  - Longitud: 2-100 caracteres
  - No permite caracteres especiales
- **Email**:
  - Formato válido
  - Único en el sistema 
  - Máximo 254 caracteres
- **Teléfono**: 
  - Mínimo 7 dígitos
  - Solo números
- **Edad**:
  - Rango: 1-150 
  - Debe coincidir con fecha de nacimiento
- **Fecha de Nacimiento**:
  - Formato ISO 8601
  - No puede ser futura
  - Debe coincidir con edad

## Testing y Desarrollo

### Ejecutar Pruebas
```bash
# Ejecutar todas las pruebas
go test ./... -v

# Ejecutar pruebas de integración  
go test ./integration/... -v
```

### Postman Collection

Para probar la API, puedes importar la colección de Postman usando el siguiente ID:

**ID**: cac1bf80-1783-494f-8478-96aa2ad322bf

### Solución de Problemas Comunes

1. **Servicios no inician**:
   - Verificar puertos disponibles
   - Revisar logs  
   - Comprobar Docker

2. **Fallo en construcción**: 
   - Limpiar imágenes Docker  
   - Actualizar código

Para problemas adicionales, contactar al equipo de desarrollo.

## Despliegue en AWS Lambda

Para desplegar la aplicación en AWS Lambda:

1. Construir el paquete Lambda:
```bash  
make build-lambda
```

## Documentación con Swagger

1. Asegúrate de que el servidor esté en ejecución.

2. Abre tu navegador web y ve a la siguiente URL:
   ```
   http://localhost:8100/swagger  
   ```

3. Verás la interfaz de usuario de Swagger, donde podrás explorar y probar los diferentes endpoints de la API.

## Generación de archivos de Swagger

Para generar los archivos de Swagger, se proporciona un script de Bash llamado `swagger`. Sigue estos pasos para ejecutar el script:

1. Abre una terminal en el directorio raíz del proyecto. 

2. Asegúrate de que el script `swagger` tenga permisos de ejecución. Si no los tiene, puedes otorgarlos con el siguiente comando:
   ```
   chmod +x ./scripts/swagger
   ```

3. Ejecuta el script con el siguiente comando:
   ``` 
   ./scripts/swagger
   ```

   Esto generará los archivos necesarios para la documentación de Swagger.

Una vez que hayas seguido estos pasos, podrás acceder a la documentación de Swagger en `http://localhost:8100/swagger` y explorar los diferentes endpoints disponibles en la API REST.

## Challenge

Descripción del Desafío (con Docker y preparado para Lambda + KPI de Clientes):

El reto consiste en desarrollar una API en Golang que permita gestionar clientes, usando Docker para contenerizar la aplicación y dejando todo listo para que la API pueda ser fácilmente desplegada en AWS Lambda. Además, se debe permitir la carga de edad y fecha de nacimiento de los clientes, y se debe implementar un endpoint que devuelva KPI de los clientes, como el promedio de edad y la desviación estándar.

Requisitos del Desafío:
   - Stack Tecnológico:
   - Lenguaje: Golang.  
   - Framework HTTP: Echo o Gin.
   - Base de datos: SQLite (o una base de datos en memoria para simplificar).

Contenerización: Docker.
Gestión de dependencias: Utilizar Go Modules.  
Documentación de la API: Implementar Swagger para la documentación automática de la API.
Preparado para AWS Lambda: El código debe estar estructurado para facilitar su despliegue en AWS Lambda o entornos similares.

Endpoints Requeridos:
   - GET /clients: Obtiene la lista de todos los clientes desde la base de datos en SQLite o en memoria.
   - GET /clients/{id}: Obtiene los detalles de un cliente específico.
   - POST /clients: Crea un nuevo cliente. El cuerpo del request debe incluir el nombre, apellido, email, número de teléfono, edad y fecha de nacimiento.   
   - PUT /clients/{id}: Actualiza un cliente existente.
   - DELETE /clients/{id}: Elimina un cliente específico.
   - GET /clients/kpi: Devuelve KPI de los clientes, tales como:
       Promedio de edad. 
       Desviación estándar de la edad.

Validaciones de Entrada: 
   - Los campos name, last_name, email, edad y fecha de nacimiento son obligatorios para crear un cliente.
   - Validar que el email tenga un formato correcto. 
   - El número de teléfono debe ser numérico y tener un mínimo de 7 dígitos.
   - La fecha de nacimiento debe ser válida y coherente con la edad provista.
   - Manejo de errores claros y consistentes para casos como cliente no encontrado, datos inválidos, etc.

Persistencia:
Los clientes deben almacenarse localmente en una base de datos SQLite o en una estructura en memoria, permitiendo un despliegue rápido.

Contenerización con Docker:
   - Crear un Dockerfile para la aplicación Golang, con una imagen base ligera (por ejemplo, golang:alpine).  
   - El contenedor debe exponer el puerto en el que la API estará escuchando.
   - Se debe poder ejecutar la aplicación localmente usando Docker para pruebas y desarrollo.

Preparación para AWS Lambda: 
 - El código debe estar preparado para ser empaquetado y desplegado en un entorno como AWS Lambda. Se debe incluir:
    - Handler adaptado para Lambda (lambda.HandlerFunc) o estructura similar.
    - Instrucciones claras para empaquetar la aplicación en un archivo ZIP compatible con AWS Lambda.

KPI de Clientes: 
Crear un endpoint GET /clients/kpi que calcule y devuelva los siguientes indicadores clave:
      - Promedio de edad: El cálculo del promedio de edad entre todos los clientes registrados.
  - Desviación estándar de edad: Calcular la variación de las edades respecto al promedio. 
  - Los KPI deben calcularse en tiempo real basados en los clientes almacenados.

Testing:
   - Escribir pruebas unitarias para los handlers (endpoints). 
   - Crear pruebas de integración para asegurarse de que la API funciona correctamente.

Estructura del Proyecto:
/handlers: Los handlers o controladores para manejar las peticiones HTTP.
/models: Las estructuras de datos y lógica de negocio para los clientes.  
/middleware: Cualquier middleware necesario, como autenticación.
/database: Configuración de la base de datos SQLite o lógica de persistencia en memoria. 
/docker: Archivos relacionados con Docker, como el Dockerfile.
/lambda: Lógica o configuraciones necesarias para empaquetar y ejecutar la aplicación en AWS Lambda.

Requisitos Técnicos Específicos:
  - Escalabilidad: El código debe estar modularizado y preparado para escalar, tanto a nivel local como en entornos serverless.
  - Optimización: Se valorará la optimización para un entorno serverless, como minimizar el tamaño del contenedor y reducir tiempos de arranque.

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
```

Estos cambios deberían mejorar la estructura y legibilidad del README. Recuerda revisar los enlaces internos después de realizar estos cambios para asegurarte de que funcionan correctamente. Si tienes alguna otra pregunta, no dudes en preguntar.