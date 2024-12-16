Una **función closure** (o simplemente "closure") en programación es una función que "cierra" sobre su entorno léxico. Esto significa que una función closure puede capturar y recordar las variables del entorno donde fue definida, incluso después de que ese entorno haya finalizado su ejecución. En otras palabras, una closure retiene su acceso a las variables y parámetros del contexto en el que se creó.

Las closures son útiles para crear funciones con estado, para la encapsulación y para crear funciones de orden superior. A continuación, se muestra una explicación más detallada con ejemplos en Go.

### Ejemplo de Closure en Go

#### Ejemplo 1: Contador Simple

En este ejemplo, se crea una función que devuelve una función closure. Esta función closure incrementa y devuelve un contador cada vez que se llama.

```go
package main

import "fmt"

func main() {
    // Llama a la función closure
    counter := makeCounter()

    fmt.Println(counter()) // Imprime: 1
    fmt.Println(counter()) // Imprime: 2
    fmt.Println(counter()) // Imprime: 3
}

func makeCounter() func() int {
    count := 0 // Variable capturada por la closure

    // Retorna una función closure
    return func() int {
        count++
        return count
    }
}
```

#### Explicación:
1. La función `makeCounter` define una variable `count` inicializada en 0.
2. `makeCounter` devuelve una función closure que incrementa `count` y luego la devuelve.
3. Cada vez que se llama a `counter()`, la función closure recuerda el valor de `count` y lo incrementa.

#### Ejemplo 2: Filtrado de Lista

En este ejemplo, se utiliza una closure para filtrar una lista de enteros.

```go
package main

import "fmt"

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Closure que filtra números pares
    even := filter(numbers, func(n int) bool {
        return n%2 == 0
    })

    fmt.Println(even) // Imprime: [2 4 6 8 10]
}

func filter(numbers []int, test func(int) bool) []int {
    var result []int
    for _, n := range numbers {
        if test(n) {
            result = append(result, n)
        }
    }
    return result
}
```

#### Explicación:
1. La función `filter` toma una lista de enteros y una función `test` que define el criterio de filtrado.
2. `filter` aplica la función `test` a cada elemento de la lista y devuelve una nueva lista con los elementos que cumplen el criterio.
3. Se pasa una función closure que define que un número es par (`n%2 == 0`).

### Características de las Closures

- **Capturan el Entorno Léxico**: Una closure puede capturar y utilizar variables de su entorno, lo que permite crear funciones que recuerdan el estado de su contexto de creación.
- **Encapsulación**: Las closures permiten encapsular lógica y estado dentro de una función, lo que puede ser útil para crear funciones con estado o comportamientos específicos.
- **Funciones de Orden Superior**: Las closures son ampliamente utilizadas en funciones de orden superior, donde las funciones pueden ser pasadas como argumentos y retornadas desde otras funciones.

### Beneficios de las Closures

- **Flexibilidad y Reutilización**: Las closures permiten crear funciones flexibles y reutilizables que pueden comportarse de manera diferente dependiendo de las variables capturadas.
- **Modularidad**: Facilitan la creación de código modular y limpio, ya que las closures pueden encapsular detalles de implementación y estado interno.
- **Mantenimiento del Estado**: Permiten mantener el estado a través de múltiples llamadas a la función closure sin necesidad de variables globales.

En resumen, las closures son una característica poderosa en muchos lenguajes de programación, incluyendo Go, que permiten capturar y recordar el entorno léxico, facilitando la creación de funciones con estado y la encapsulación de lógica.