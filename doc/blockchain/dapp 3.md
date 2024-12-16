# Aplicación de Encuestas con Blockchain, Go y zkML

Aplicación de encuestas que utiliza blockchain, Go, y zkML para asegurar que los votos sean inviolables e inalterables.

## Estructura del Proyecto

1. **Blockchain**:
   - **Ethereum**: Almacena las encuestas y los votos de manera descentralizada e inmutable.
   - **Contratos Inteligentes**: Implementados en Solidity para gestionar las encuestas y los votos.

2. **Backend en Go**:
   - **API**: Actúa como interfaz entre el frontend y los contratos inteligentes en Ethereum.
   - **Generación de Pruebas zk-SNARK**: Verifica la validez de los votos sin revelar información privada.

3. **zkML (Zero-Knowledge Machine Learning)**:
   - **Pruebas zk-SNARK**: Aseguran que los votos sean válidos y cumplen con los requisitos de privacidad.

4. **Frontend**:
   - **Interfaz de Usuario**: Permite a los usuarios crear encuestas, votar y ver resultados.
   - **Interacción con Blockchain**: Utiliza Web3.js o ethers.js para conectar con los contratos inteligentes.

## Pasos para Implementar el Proyecto

### 1. Desarrollar Contratos Inteligentes en Solidity

Diseñar y codificar contratos inteligentes que gestionen la creación de encuestas y la emisión de votos. Estos contratos incluyen funcionalidades para:
- Crear una nueva encuesta.
- Emitir un voto en una encuesta existente.
- Contar los votos y mostrar los resultados de manera transparente.

### 2. Desplegar Contratos Inteligentes

Desplegar los contratos inteligentes en la red Ethereum. Esto implica cargar el código del contrato en la blockchain y asegurar su correcta ejecución. Es importante realizar pruebas exhaustivas en una red de prueba (testnet) antes del despliegue en la red principal (mainnet).

### 3. Desarrollar el Backend en Go

Implementar una API RESTful en Go que sirva como intermediaria entre el frontend y los contratos inteligentes en Ethereum. La API gestiona:
- Solicitudes para crear nuevas encuestas.
- Emisión de votos.
- Consulta de resultados.

### 4. Implementar zk-SNARKs

Integrar zk-SNARKs para asegurar la privacidad y validez de los votos. Esto incluye:
- Desarrollar circuitos zk-SNARK para verificar los votos sin revelar detalles sensibles.
- Implementar funciones en el backend para generar y verificar estas pruebas antes de registrar un voto en la blockchain.

### 5. Desarrollar el Frontend

Crear una interfaz de usuario accesible y amigable que permita a los usuarios:
- Crear nuevas encuestas.
- Emitir votos de manera intuitiva y segura.
- Consultar resultados de encuestas en tiempo real.

La interfaz utiliza bibliotecas como Web3.js o ethers.js para interactuar con los contratos inteligentes y permitir la conexión a la blockchain desde el navegador del usuario.

### 6. Pruebas y Seguridad

Realizar pruebas exhaustivas de todos los componentes del sistema para asegurar su correcto funcionamiento y seguridad. Esto incluye:
- Pruebas unitarias y de integración para los contratos inteligentes.
- Pruebas de carga y estrés para la API y el frontend.
- Auditorías de seguridad para identificar y mitigar posibles vulnerabilidades en el sistema.
