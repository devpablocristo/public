La **multi-tenencia (multi-tenancy)** es un modelo de arquitectura de software en el cual una única instancia de un sistema de software sirve a múltiples clientes (denominados "inquilinos" o "tenants"). Cada inquilino es un grupo distinto de usuarios que comparte acceso a la instancia de software, pero sus datos y configuraciones están aislados y protegidos de otros inquilinos. Este modelo es común en aplicaciones de Software como Servicio (SaaS).

### Características Clave de la Multi-Tenencia

1. **Aislamiento de Datos**:
   - Aunque múltiples inquilinos comparten la misma instancia de software, sus datos están completamente separados. Esto asegura que un inquilino no pueda acceder a los datos de otro.

2. **Escalabilidad**:
   - Permite que una aplicación escale eficientemente ya que múltiples inquilinos pueden ser servidos desde una única infraestructura.

3. **Mantenimiento Simplificado**:
   - Dado que sólo hay una instancia del software que mantener, las actualizaciones y el mantenimiento son más sencillos y eficientes en comparación con mantener instancias separadas para cada cliente.

4. **Personalización**:
   - Los inquilinos pueden personalizar aspectos de la aplicación para satisfacer sus necesidades específicas sin afectar a otros inquilinos. Esto incluye configuraciones, interfaces de usuario y a veces incluso funcionalidades.

### Modelos de Multi-Tenencia

1. **Modelo con una Base de Datos por Inquilino**:
   - Cada inquilino tiene su propia base de datos, lo que proporciona un alto grado de aislamiento de datos, pero puede ser más complejo de gestionar a medida que aumenta el número de inquilinos.

2. **Modelo con Esquemas Separados**:
   - Todos los inquilinos comparten una única base de datos, pero cada uno tiene su propio esquema dentro de la base de datos. Esto equilibra el aislamiento de datos con la facilidad de gestión.

3. **Modelo con Tablas Compartidas**:
   - Todos los inquilinos comparten las mismas tablas dentro de una base de datos, y se utiliza una columna para identificar el inquilino al que pertenecen los datos. Este modelo es el más eficiente en términos de recursos, pero también el más complejo de implementar en cuanto a la seguridad y el aislamiento de datos.

### Ventajas de la Multi-Tenencia

- **Reducción de Costes**:
  - Al compartir una única instancia de software y hardware entre múltiples inquilinos, los costos de operación y mantenimiento se reducen significativamente.

- **Eficiencia Operativa**:
  - Las actualizaciones y mejoras se despliegan una sola vez y están disponibles para todos los inquilinos, simplificando el ciclo de desarrollo y despliegue.

- **Mejor Utilización de Recursos**:
  - El uso de recursos se optimiza al servir a múltiples inquilinos desde una única infraestructura, lo que permite una mejor escalabilidad y gestión de recursos.

### Desafíos de la Multi-Tenencia

- **Seguridad de Datos**:
  - Garantizar que los datos de un inquilino no sean accesibles para otros es crucial y puede ser complejo de implementar.

- **Personalización**:
  - Balancear la personalización de la aplicación para satisfacer las necesidades de diferentes inquilinos sin afectar la integridad del sistema compartido.

- **Gestión de Rendimiento**:
  - Asegurar que la aplicación funcione de manera eficiente para todos los inquilinos, incluso cuando la carga de trabajo varía entre ellos.

En resumen, la multi-tenencia es una estrategia eficaz para maximizar la eficiencia y reducir costos en aplicaciones SaaS, aunque presenta desafíos significativos en términos de seguridad y gestión de datos.

Sí, puedes usar tanto multi-tenancy como la arquitectura hexagonal (hex arch) en un proyecto. De hecho, combinarlas puede ser una excelente manera de estructurar una aplicación robusta y escalable. A continuación, te explico cómo puedes integrarlas:

### Integración de Multi-Tenancy y Arquitectura Hexagonal

**Multi-Tenancy** se centra en gestionar múltiples inquilinos dentro de una misma aplicación, asegurando el aislamiento de datos y personalización por inquilino. **La Arquitectura Hexagonal** se enfoca en la separación de preocupaciones dentro de la aplicación, dividiendo la lógica de negocio, interfaces y la infraestructura.

### Cómo se pueden combinar:

1. **Lógica de Negocio (Núcleo) y Multi-Tenancy**:
    - **Lógica de Negocio Independiente**: Mantén la lógica de negocio independiente de los detalles de implementación específicos del inquilino. Esto se logra definiendo puertos que abstraen las operaciones.
    - **Identificación de Inquilinos**: Incluye la identificación del inquilino como parte del modelo de dominio cuando sea necesario. Por ejemplo, al definir una entidad como `Order`, puede incluir un campo `TenantID`.

2. **Puertos y Adaptadores**:
    - **Puertos (Interfaces)**: Define puertos que soporten operaciones multi-tenancy, como `LoadTenantData(tenantID string)`.
    - **Adaptadores**: Implementa adaptadores específicos para manejar el aislamiento de datos y las configuraciones personalizadas para cada inquilino.
        - **Adaptadores de Persistencia**: Puedes tener adaptadores que manejen conexiones a diferentes bases de datos o esquemas basados en el inquilino.
        - **Adaptadores de Interfaz de Usuario**: Implementa adaptadores que manejen interfaces de usuario personalizadas por inquilino si es necesario.

3. **Configuración y Enrutamiento**:
    - **Middlewares**: Utiliza middlewares para identificar al inquilino en las solicitudes entrantes y establecer el contexto adecuado.
    - **Enrutamiento**: Configura el enrutamiento de manera que soporte las rutas específicas por inquilino.

### Ejemplo de Implementación:

#### Definición de Entidades y Repositorios (Núcleo):

```go
package item

import (
    "time"
)

type Item struct {
    ID        int       `json:"id"`
    TenantID  string    `json:"tenant_id"` // Identificador del inquilino
    Code      string    `json:"code"`
    Title     string    `json:"title"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ItemRepositoryPort interface {
    SaveItem(item *Item) error
    ListItems(tenantID string) ([]Item, error) // Método específico de multi-tenancy
}
```

#### Implementación de Repositorios (Adaptadores):

```go
package item

import (
    "database/sql"
)

type mysqlRepository struct {
    db *sql.DB
}

func NewMySQLRepository(db *sql.DB) ItemRepositoryPort {
    return &mysqlRepository{db: db}
}

func (r *mysqlRepository) SaveItem(item *Item) error {
    query := `INSERT INTO items (tenant_id, code, title, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
    _, err := r.db.Exec(query, item.TenantID, item.Code, item.Title, item.CreatedAt, item.UpdatedAt)
    return err
}

func (r *mysqlRepository) ListItems(tenantID string) ([]Item, error) {
    query := `SELECT id, tenant_id, code, title, created_at, updated_at FROM items WHERE tenant_id = ?`
    rows, err := r.db.Query(query, tenantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []Item
    for rows.Next() {
        var item Item
        if err := rows.Scan(&item.ID, &item.TenantID, &item.Code, &item.Title, &item.CreatedAt, &item.UpdatedAt); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, nil
}
```

#### Configuración y Middleware:

```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "X-Tenant-ID header is required"})
            return
        }
        c.Set("tenantID", tenantID)
        c.Next()
    }
}
```

#### Configuración del Enrutador:

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"

    "api/rest/handlers"
    "api/internal/core"
    "api/internal/item"
    "api/internal/bootstrap/mysql"
    "api/middleware"
)

func main() {
    // Configurar MySQL
    mysqlClient, err := mysql.NewMySQLSetup()
    if err != nil {
        log.Fatalf("Could not set up MySQL: %v", err)
    }
    defer mysqlClient.Close()

    // Inicializar repositorios
    mysqlRepo := item.NewMySQLRepository(mysqlClient.DB())

    // Inicializar caso de uso
    usecase := core.NewItemUsecase(mysqlRepo)

    // Inicializar handlers
    handler := handlers.NewHandler(usecase)

    // Configurar enrutador
    router := gin.Default()
    router.Use(middleware.TenantMiddleware()) // Middleware para manejo de inquilinos

    router.POST("/items", handler.SaveItem)
    router.GET("/items", handler.ListItems)

    // Iniciar servidor
    log.Println("Server started at http://localhost:8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

### Resumen:

- **Multi-Tenancy** se encarga de manejar múltiples inquilinos, asegurando el aislamiento de datos y la personalización de cada uno.
- **Arquitectura Hexagonal** proporciona una separación clara de la lógica de negocio, interfaces y adaptadores, facilitando la testabilidad y flexibilidad.

Al combinar ambas, se obtiene una aplicación robusta, escalable y mantenible, que puede servir a múltiples clientes de manera eficiente y segura.