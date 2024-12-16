Paso 1 ok
Paso 2 ok
Paso 3 ok
Paso 4 ok
Paso 5 ok

REVISAR
***********************************************************************
Paso 6: Almacenamiento del Token de Verificación
MS: auth
Acción: Guardar el token generado.
Descripción:
auth almacena el token de verificación en Redis, asociado al userUUID y con un tiempo de expiración definido.
Protocolo de Comunicación: Interno (auth ⇔ Redis)
***********************************************************************

Paso 7: Envío del Correo de Verificación
MS Involucrados:
Publicador: users
Consumidor: mailing
Acción: users envía un mensaje a la cola para mailing.
Descripción:
Si se requiere verificación adicional del correo:
users coloca un mensaje en la cola (AWS SQS) con el correo electrónico y el token de verificación.

mailing consume el mensaje, construye el correo electrónico y lo envía al usuario utilizando AWS SES.

Protocolo de Comunicación:
Mensajería Asíncrona (users ⇔ AWS SQS ⇔ mailing)
API de AWS SES para enviar el correo electrónico.

hasta aqui - falta seguir. 