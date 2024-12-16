### Documentación de Google Wire

#### Introducción

`Google Wire` es una herramienta de inyección de dependencias para Go que permite a los desarrolladores conectar componentes de manera declarativa. Utiliza el análisis estático para generar el código necesario para la inyección de dependencias, reduciendo la necesidad de escribir código repetitivo y mejorando la mantenibilidad del proyecto.

#### Instalación

Para instalar `Google Wire`, ejecuta el siguiente comando:

```sh
go install github.com/google/wire/cmd/wire@latest
```

#### Conceptos Básicos

1. **Proveedores (Providers)**:
   Un proveedor es una función que crea y configura un tipo específico. Los proveedores pueden depender de otros tipos que también deben ser inyectados. Se definen utilizando funciones que retornan instancias de los tipos requeridos.

   ```go
   func provideFoo() Foo {
       return Foo{}
   }
   ```

2. **Sets de Proveedores (Provider Sets)**:
   Los sets de proveedores son colecciones de proveedores que `Wire` utiliza para entender cómo construir un grafo de dependencias. Se crean usando `wire.NewSet` y pueden incluir funciones de proveedor, otros sets de proveedores y valores ya existentes.

   ```go
   var mySet = wire.NewSet(provideFoo, provideBar)
   ```

3. **Inyectores (Injectors)**:
   Un inyector es una función especial que `Wire` genera. Esta función inicializa todas las dependencias necesarias y las conecta. Los inyectores se definen como funciones que utilizan `wire.Build` para declarar qué proveedores y sets de proveedores se deben usar para resolver las dependencias.

   ```go
   func InitializeFoo() (*Foo, error) {
       wire.Build(mySet)
       return &Foo{}, nil
   }
   ```

4. **Build Tags**:
   `Wire` utiliza una etiqueta de compilación `// +build wireinject` para evitar que el archivo se compile normalmente. Esto permite a `Wire` generar el archivo con el código necesario para la inyección de dependencias.

#### Pasos para Usar Wire

1. **Definir Proveedores**:
   Define funciones proveedoras que crean y configuran instancias de los tipos que necesitas. Estas funciones deben estar en paquetes donde se puedan acceder a todas las dependencias necesarias.

   ```go
   func provideFoo() Foo {
       return Foo{}
   }

   func provideBar() Bar {
       return Bar{}
   }
   ```

2. **Crear Sets de Proveedores**:
   Crea sets de proveedores utilizando `wire.NewSet`, agrupando todas las funciones proveedoras y otros sets de proveedores necesarios para inicializar los tipos requeridos.

   ```go
   var mySet = wire.NewSet(provideFoo, provideBar)
   ```

3. **Definir Inyectores**:
   Define funciones inyectoras que utilizan `wire.Build` para especificar qué sets de proveedores deben usarse para resolver las dependencias. Estas funciones deben estar marcadas con la etiqueta `// +build wireinject`.

   ```go
   // +build wireinject

   package mypackage

   func InitializeFoo() (*Foo, error) {
       wire.Build(mySet)
       return &Foo{}, nil
   }
   ```

4. **Generar el Código**:
   Ejecuta el comando `wire` en el directorio del paquete para generar el código de inyección de dependencias. `Wire` creará un archivo con el sufijo `_wire_gen.go` que contiene el código generado.

   ```sh
   wire
   ```

   ```sh
   wire ./...
   ```

   El primer comando genera archivos de inyección de dependencias solo en el directorio actual, mientras que el segundo comando busca recursivamente en el directorio actual y en todos los subdirectorios, generando archivos de inyección de dependencias para todos los archivos `wire.go` que encuentra.

#### Ubicación de los Archivos Generados

Los archivos generados por `wire` se ubican en los mismos directorios que contienen los archivos `wire.go`. Aquí hay un ejemplo de cómo se organiza un proyecto y dónde se generan los archivos:

##### Estructura del Proyecto

```plaintext
myproject/
|-- cmd/
|   |-- api/
|   |   |-- wire.go
|-- internal/
|   |-- core/
|   |   |-- wire.go
|   |-- platform/
|       |-- config/
|       |   |-- wire.go
```

##### Archivos Generados por `wire`

Después de ejecutar `wire ./...`, la estructura del proyecto incluirá los archivos generados:

```plaintext
myproject/
|-- cmd/
|   |-- api/
|   |   |-- wire.go
|   |   |-- wire_gen.go
|-- internal/
|   |-- core/
|   |   |-- wire.go
|   |   |-- wire_gen.go
|   |-- platform/
|       |-- config/
|       |   |-- wire.go
|       |   |-- wire_gen.go
```

#### Usar el Código Generado

Usa la función inyectora generada en tu aplicación para inicializar y utilizar las dependencias.

```go
package main

func main() {
    foo, err := InitializeFoo()
    if err != nil {
        log.Fatal(err)
    }
    // Usa foo...
}
```

#### Beneficios de Usar Wire

1. **Mantenimiento Mejorado**:
   Facilita el mantenimiento del código al centralizar la configuración y conexión de dependencias, lo que reduce la complejidad y la repetición.

2. **Reducción de Errores**:
   Disminuye la probabilidad de errores de configuración de dependencias al generar automáticamente el código necesario para conectar componentes.

3. **Modularidad y Reutilización**:
   Los proveedores y sets de proveedores son reutilizables en diferentes contextos, mejorando la modularidad y la capacidad de prueba del código.

#### Consideraciones Finales

- **Verificación Estática**:
  `Wire` realiza verificación estática del grafo de dependencias en tiempo de compilación, lo que ayuda a detectar errores de configuración de dependencias temprano en el proceso de desarrollo.

- **Compatibilidad**:
  `Wire` está diseñado para ser compatible con las prácticas de codificación idiomáticas de Go, utilizando tipos y funciones en lugar de reflexión o generación dinámica de código en tiempo de ejecución.

- **Documentación Adicional**:
  Para más detalles, puedes consultar la [documentación oficial de `Wire`](https://github.com/google/wire).

Con estos conceptos y pasos, puedes empezar a usar `Google Wire` para gestionar las dependencias en tus proyectos Go, mejorando la organización y mantenibilidad del código.