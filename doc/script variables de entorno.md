Al modificar las variables de entorno, actualmente es necesario **detener y reiniciar los contenedores** para que los cambios se apliquen, sin necesidad de reconstruir las imágenes. 

Una solución sencilla, pero posiblemente incompleta, es utilizar el paquete `godotenv` para cargar las variables de entorno en la aplicación Go. Esta estrategia funcionaría bien para la propia aplicación, ya que `godotenv` permite recargar las variables de entorno sin reiniciar el contenedor. Sin embargo, presenta una limitación significativa: **si existen otros servicios dentro de los contenedores que también dependen de estas variables de entorno, no se actualizarían automáticamente**. Estos servicios no serían conscientes de los cambios realizados en los archivos `.env` y continuarían operando con las configuraciones antiguas, lo que podría llevar a inconsistencias o comportamientos inesperados en el sistema.

### **Resumen de la Situación:**

1. **Actualización de Variables de Entorno:**
   - **Necesidad:** Al modificar las variables de entorno, es necesario reiniciar los contenedores para que los cambios surtan efecto.
   - **Solución Actual:** Uso de `godotenv` para cargar variables en la aplicación Go, evitando así la necesidad de reconstruir la imagen.

2. **Limitaciones de la Solución con `godotenv`:**
   - **Alcance:** Solo aplica a la aplicación Go.
   - **Problema:** Otros servicios dentro de los mismos contenedores que dependen de las mismas variables de entorno no detectan los cambios y continúan operando con las configuraciones anteriores.

### **Consideraciones Adicionales:**

- **Consistencia entre Servicios:** En entornos donde múltiples servicios dependen de un conjunto común de variables de entorno, es crucial asegurar que todos ellos reciban y apliquen las actualizaciones de manera coherente.
  
- **Automatización y Eficiencia:** Reiniciar manualmente los contenedores cada vez que se actualicen las variables de entorno puede ser propenso a errores y consumir tiempo, especialmente en entornos de desarrollo o despliegue continuo.