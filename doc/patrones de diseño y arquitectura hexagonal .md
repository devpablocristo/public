# Guía de Diseño y Arquitectura de Software en Go: Patrones de Diseño y Arquitectura Hexagonal

## Introducción
En esta guía discutiremos los patrones de diseño y arquitectura en el desarrollo de software, centrándonos en la implementación en Go con el ORM GORM. Veremos patrones clave como Data Transfer Object (DTO), Data Access Object (DAO), y Entity, así como el patrón Repository y cómo estos se integran en la Arquitectura Hexagonal.

## GORM y gorm.Model
GORM es un ORM (Object-Relational Mapper) popular en Go que proporciona una API de alto nivel para realizar operaciones de base de datos. En GORM, `gorm.Model` es una estructura básica que incluye algunos campos comunes: ID, CreatedAt, UpdatedAt, DeletedAt. Se puede incrustar en tus propias estructuras para añadir estos campos automáticamente.

```go
type User struct {
    gorm.Model
    Name  string
    Email string
}
```

## DTO, Entity, DAO
Estos tres patrones se utilizan comúnmente en el diseño de software.

- **DTO (Data Transfer Object)**: Se utiliza para transferir datos entre procesos o componentes de la aplicación. Suele ser una estructura simple que sólo agrupa datos.

```go
type UserDTO struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

- **Entity**: Es una representación de un objeto de dominio en la aplicación. Contiene tanto datos como comportamiento.

```go
type User struct {
    ID    uint
    Name  string
    Email string
}
```

- **DAO (Data Access Object)**: Proporciona una interfaz para el acceso a los datos a nivel de una entidad específica. Se ocupa principalmente de las operaciones CRUD (Crear, Leer, Actualizar, Eliminar).

```go
type UserDAO struct {
    db *gorm.DB
}

func (dao *UserDAO) FindByID(id uint) (*User, error) { /*...*/ }
func (dao *UserDAO) Save(user *User) error { /*...*/ }
```

En algunos casos, las estructuras de la entidad de dominio y de la entidad de persistencia pueden variar en función de las necesidades de la aplicación, pero en otros casos pueden ser la misma estructura. 

## Uso de Entidades Personalizadas en Presenters y Persistencia

En ciertos casos, es beneficioso definir entidades personalizadas para los presenters y la persistencia.

**Presenter con Entidades Personalizadas**: Un Presenter puede tener una entidad personalizada. Es parte del patrón de diseño Model-View-Presenter (MVP), donde se toman los datos del modelo y se preparan para la vista.

```go
type User struct {
    ID    uint
    Name  string
    Email string
    Password string
    Orders []Order
}

type UserPresenter struct {}

type UserDTO struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    OrderCount int `json:"order_count"`
}

func (p *UserPresenter) Present(user *User) *UserDTO {
    return &UserDTO{
        ID:    strconv.FormatUint(uint64(user.ID), 10),
        Name:  user.Name,
        Email: user.Email,
        OrderCount: len(user.Orders),
    }
}
```

**Persistencia con Entidades Personal

izadas**: Para el manejo de la persistencia de datos, es común utilizar una entidad personalizada que representa la estructura de la tabla de la base de datos.

```go
type UserModel struct {
    gorm.Model
    Name  string
    Email string
}
```

## Patrón Repository vs Patrón DAO
Ambos patrones proporcionan abstracción sobre las operaciones de acceso a los datos, pero se diferencian en sus enfoques y usos típicos.

**DAO (Data Access Object)**: Proporciona una abstracción de cualquier tipo de operación de persistencia y se asocia con operaciones de nivel de tabla específicas en una base de datos SQL.

```go
type UserDAO struct {
    db *gorm.DB
}

func (dao *UserDAO) Insert(user *User) error { /*...*/ }
func (dao *UserDAO) Update(user *User) error { /*...*/ }
```

**Repository**: Añade una capa de abstracción sobre las operaciones de almacenamiento y recuperación de objetos de dominio y se adhiere más a un estilo orientado a objetos de manipulación de entidades.

```go
type UserRepository interface {
    Save(user *User) error
}

type GormUserRepository struct {
    db *gorm.DB
}

func (repo *GormUserRepository) Save(user *User) error { /*...*/ }
```

En el caso de la arquitectura hexagonal, se tiende a preferir el patrón Repository porque ofrece una mayor desacoplamiento entre la lógica de dominio y la infraestructura de persistencia.

## Conclusión
En esta guía, hemos cubierto varios patrones de diseño y arquitectura de software comunes y cómo pueden implementarse en Go usando GORM. Los patrones de diseño, como DTO, DAO, y Repository, así como el patrón Presenter, son herramientas útiles para mantener tu código limpio, organizado, y fácil de entender. Finalmente, discutimos cómo estos patrones se integran en la arquitectura hexagonal, proporcionando una sólida arquitectura de aplicación que es fácil de mantener y evolucionar.
