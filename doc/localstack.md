para crear en una cola, entrar en el contenedor de localstack y correr por la consola:

aws --endpoint-url=http://localhost:4566 --region us-east-1 sqs create-queue --queue-name users_verification_email

para recibir el mensaje de la cola

aws --endpoint-url=http://localhost:4566 --region us-east-1 sqs receive-message --queue-url http://localhost:4566/000000000000/users_verification_email

