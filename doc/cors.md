### CORS (Cross-Origin Resource Sharing)

**CORS** es una característica de seguridad del navegador que permite o restringe las solicitudes realizadas a un dominio diferente (origen) desde el que se sirvió el recurso. Es esencialmente una política de seguridad que protege a los navegadores de realizar solicitudes a dominios que no son de confianza.

#### ¿Qué es?
- **Definición**: CORS es un mecanismo que utiliza cabeceras HTTP adicionales para permitir que un servidor indique cualquier otro origen (dominio, esquema o puerto) desde el cual un navegador debe permitir la carga de recursos.
- **Objetivo**: Proteger la seguridad del usuario previniendo ataques de Cross-Site Scripting (XSS) y Cross-Site Request Forgery (CSRF).

#### ¿Cómo Funciona?
1. **Solicitud Previa (Preflight Request)**: Para ciertas solicitudes HTTP (como POST, PUT, DELETE), el navegador primero envía una solicitud preliminar (preflight) usando el método OPTIONS para verificar si el servidor permite la solicitud real.
2. **Cabeceras Importantes**:
   - **Access-Control-Allow-Origin**: Especifica qué orígenes están permitidos.
   - **Access-Control-Allow-Methods**: Especifica los métodos HTTP permitidos (GET, POST, etc.).
   - **Access-Control-Allow-Headers**: Indica qué cabeceras pueden ser usadas durante la solicitud real.
   - **Access-Control-Allow-Credentials**: Indica si se permite el envío de credenciales (cookies, encabezados de autenticación, etc.).

#### Ejemplo de Cabeceras CORS:
```http
Access-Control-Allow-Origin: https://example.com
Access-Control-Allow-Methods: GET, POST, PUT
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Allow-Credentials: true
```

En resumen, **CORS** es una política de seguridad para controlar el acceso a recursos desde diferentes orígenes.