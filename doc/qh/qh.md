# Descripción General
Cada uno de estos microservicios es parte de una red social centrada en eventos, diseñada para ofrecer una experiencia interactiva y enriquecedora a los usuarios, permitiéndoles descubrir, participar y evaluar eventos dentro de su comunidad y más allá.

Para crear un microservicio usar como modelo: [analitics](https://github.com/devpablocristo/events-sn/analitics)

Definir las rutas con "/api/v1/nombre", ejemplo:

```go
v1 := r.Group("/api/v1/analytics")
{
  v1.POST("/create-report", handler.CreateReport)
}
```

## 1. Microservicio de Gestión de Eventos (Lu)
- **Endpoint:** `/events`
- **Descripción:** Responsable de todas las operaciones relacionadas con la gestión de eventos, incluyendo la creación, modificación, eliminación y consulta de eventos. Administra entidades principales del dominio como Evento, Ubicación y Categoría.
- **Entidades Principales:** Evento, Ubicación, Categoría.
- **Casos de Uso:** CrearEvento, ModificarEvento, EliminarEvento, ConsultarEvento.
- **Interacciones:**
  - Con Notificaciones: Informa sobre cambios en los eventos.
  - Con Historial y Archivos: Envía datos de eventos concluidos para su almacenamiento y análisis.
- **Puertos:**
  - Entrada: API REST que recibe solicitudes.
  - Salida: Interfaces para la persistencia de eventos y comunicación con otros servicios.
- **Adaptadores:**
  - Primarios: Controladores REST.
  - Secundarios: Implementaciones de la base de datos PostgreSQL y producción de mensajes para notificaciones.
- **Código:** `event.go`

```go
package event

type Event struct {
    ID          string
    Title       string
    Description string
    Location    string
    StartTime   time.Time
    EndTime     time.Time
    Category    string
}
```

## 2. Microservicio de Búsqueda y Recomendaciones (Mau)
- **Endpoint:** `/search`
- **Descripción:** Facilita la búsqueda de eventos basada en diversos criterios y proporciona recomendaciones personalizadas a los usuarios. Utiliza tecnologías de búsqueda avanzada para ofrecer resultados eficientes y relevantes.
- **Entidades Principales:** Búsqueda, Recomendación.
- **Casos de Uso:** RealizarBúsqueda, GenerarRecomendaciones.
- **Interacciones:**
  - Con Gestión de Eventos: Recupera detalles de eventos para procesar búsquedas y generar recomendaciones.
  - Con Usuarios y Autenticación: Obtiene información del perfil del usuario para personalizar las recomendaciones de eventos.
- **Puertos:**
  - Entrada: API REST.
  - Salida: Interfaces para obtener datos de eventos y de usuarios.
- **Adaptadores:**
  - Primarios: Controladores REST.
  - Secundarios: Cliente HTTP para integración con otros microservicios, implementación de Elasticsearch.
- **Código:** `search.go`

```go
package search

type SearchResult struct {
    Events []event.Event
}
```

## 3. Microservicio de Calificaciones y Reseñas (Mati)
- **Endpoint:** `/reviews`
- **Descripción:** Gestiona todas las reseñas y calificaciones de eventos y organizadores de eventos. Permite a los usuarios agregar, modificar, eliminar y consultar reseñas, ayudando a otros usuarios a tomar decisiones informadas sobre los eventos a asistir.
- **Entidades Principales:** Calificación, Reseña.
- **Casos de Uso:** AñadirReseña, ModificarReseña, EliminarReseña, ConsultarReseñas.
- **Interacciones:**
  - Con Gestión de Eventos: Recupera información de eventos para asociar reseñas y calificaciones con eventos específicos.
  - Con Usuarios y Autenticación: Verifica la identidad del usuario y asegura que solo los usuarios que asistieron al evento puedan dejar una reseña.
- **Puertos:**
  - Entrada: API REST.
  - Salida: Interfaces para la persistencia y validación de usuario.
- **Adaptadores:**
  - Primarios: Controladores REST.
  - Secundarios: Implementaciones de MongoDB para almacenamiento de datos.
- **Código:** `review.go`

```go
package review

type Review struct {
    ID         string
    UserID     string
    EventID    string
    Rating     int
    Comment    string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

## 4. Microservicio de Usuarios y Autenticación (Joaco)
- **Endpoint:** `/users`
- **Descripción:** Controla la gestión de usuarios y la autenticación en la plataforma, registrando nuevos usuarios, gestionando sesiones y actualizaciones de perfil, y asegurando que las credenciales se manejen de manera segura.
- **Entidades Principales:** Usuario.
- **Casos de Uso:** RegistrarUsuario, IniciarSesión, ModificarPerfil, ValidarUsuario.
- **Interacciones:**
  - Con Notificaciones: Envía detalles de usuario para personalizar las notificaciones y gestionar preferencias de notificación.
  - Con Búsqueda y Recomendaciones: Proporciona información del perfil del usuario para personalizar las recomendaciones de eventos.
- **Puertos:**
  - Entrada: API REST.
  - Salida: Interfaces para la gestión de credenciales y datos de usuario.
- **Adaptadores:**
  - Primarios: Controladores REST.
  - Secundarios: Implementaciones de PostgreSQL para manejar datos de usuario.
- **Código:** `user.go`

```go
package user

type User struct {
    ID       string
    Username string
    Password string
    Email    string
}
```

## 5. Microservicio de Notificaciones (Gonza)
- **Endpoint:** `/notifications`
- **Descripción:** Maneja la distribución de notificaciones a los usuarios, informándoles sobre eventos relevantes como actualizaciones de eventos a los que están suscritos o nuevos eventos que podrían interesarles.
- **Entidades Principales:** Notificación.
- **Casos de Uso:** EnviarNotificación.
- **Interacciones:**
  - Con Usuarios y Autenticación: Recibe información de los usuarios para enviar notificaciones personalizadas basadas en sus intereses y preferencias.
  - Con Gestión de Eventos: Recibe alertas de eventos nuevos o modificados para notificar a los usuarios pertinentes.
- **Puertos:**
  - Entrada: API interna.
  - Salida: Interfaces para la entrega de notificaciones.
- **Adaptadores:**
  - Primarios: Consumidor de mensajes.
  - Secundarios: Implementaciones de Redis para la gestión de mensajes.
- **Código:** `notification.go`

```go
package notification

type Notification struct {
    ID      string
    UserID  string
    Message string
}
```

## 6. Microservicio de Historial y Archivos (Cande)
- **Endpoint:** `/history`
- **Descripción:** Proporciona funcionalidades para archivar eventos pasados y consultarlos, manteniendo un registro de la actividad de eventos y ofreciendo estadísticas y análisis sobre tendencias pasadas.
- **Entidades Principales:** EventoHistórico.
- **Casos de Uso:** ArchivarEvento, ConsultarHistorial.
- **Interacciones:**
  - Con Gestión de Eventos: Recibe datos de eventos pasados para archivar y proporcionar funcionalidades de consulta de eventos históricos a los usuarios y para análisis de tendencias.
  - Con Búsqueda y Recomendaciones: Provee datos históricos para influir en las recomendaciones, utilizando análisis de tendencias y popularidad pasada de eventos.
- **Puertos:**
  - Entrada: API REST.
  - Salida: Interfaces para la persistencia de eventos históricos.
- **Adaptadores:**
  - Primarios: Controladores REST.
  - Secundarios: Implementaciones de PostgreSQL y Amazon S3 para el almacenamiento.
- **Código:** `history.go`

```go
package history

type HistoricalEvent struct {
    EventID     string
    EventData   event.Event
    OccurredAt  time.Time
}
```
## 7. Microservicio de Análisis de Datos y Reportes (Pablo)
- **Endpoint:** `/analytics`
- **Descripción:** Proporciona análisis de datos y reportes detallados sobre la interacción de los usuarios con eventos, ofreciendo insights que ayudan a mejorar la organización de eventos futuros y la experiencia del usuario.
- **Entidades Principales:** Reporte, Análisis de Usuario, Métricas de Evento.
- **Casos de Uso:** GenerarReportes, AnalizarTendencias, MonitorizarActividad.
- **Interacciones:**
  - **Con Gestión de Eventos y Búsqueda y Recomendaciones:** Utiliza datos de actividad reciente para generar análisis y recomendaciones.
  - **Con Historial y Archivos:** Accede a datos archivados para análisis de tendencias a largo plazo.
- **Puertos:**
  - **Entrada:** API REST que acepta consultas de análisis.
  - **Salida:** Interfaces para la visualización y exportación de datos y reportes.
- **Adaptadores:**
  - **Primarios:** Controladores REST para recepción de solicitudes de datos.
  - **Secundarios:** Herramientas de análisis de datos como Apache Spark, implementaciones de bases de datos como Google BigQuery para el manejo de grandes volúmenes de datos.
- **Código:** `report.go`

```go
package report

import "time"

type Report struct {
	ReportID    string
	GeneratedAt time.Time
	Metrics     map[string]interface{}
}
```