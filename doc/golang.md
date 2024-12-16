Aquí te explico en detalle cómo funciona el acceso encadenado en Go:

Los casos principales para recordar son:

1. **Con Estructuras**: Puedes acceder a sus campos usando el punto.
```go
persona := ObtenerPersona().Nombre  // ✅ Válido
```

2. **Con Punteros**: Funciona igual que con estructuras directas.
```go
persona := ObtenerPersonaPuntero().Nombre  // ✅ Válido
```

3. **Con Tipos Simples**: No permiten encadenamiento.
```go
numero := ObtenerNumero().algo  // ❌ ERROR
```

4. **Con Múltiples Retornos**: Requieren asignación previa.
```go
valor, err := ObtenerDatos()  // ✅ Válido
valor := ObtenerDatos().algo  // ❌ ERROR
```

5. **Con Maps**: Usan corchetes en lugar del punto.
```go
valor := ObtenerMapa()["clave"]  // ✅ Válido
```

6. **Con Interfaces**: Puedes llamar a sus métodos definidos.
```go
texto := ObtenerLector().Leer()  // ✅ Válido
```

Puntos importantes a recordar:
- El encadenamiento solo funciona con tipos que tienen campos o métodos accesibles
- Los punteros se desreferencian automáticamente
- Múltiples retornos requieren asignación previa
- Maps usan sintaxis de corchetes []
- Puedes encadenar múltiples llamadas si cada paso retorna algo válido

Esta funcionalidad es parte fundamental de cómo Go maneja el acceso a campos y métodos, y es importante entender sus limitaciones y casos de uso correctos.
