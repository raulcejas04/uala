Requerimientos

Versiones que utilicé
golang version go1.21.4
docker version 27.3.1
docker-compose v2.18.1

Levantar kafka y redis

Entrar en los directorios ./kafka y ./redis y correr "docker-composer up -d". Si los volumenes son creados sin permisos suficientes no va a levantar, entonces hay que darle permisos 777 con chmod -R ugo+rwx y tambien chown -R usuario:usuario ./directorio.

Entrar en twitter y levantar los 2 microservicios:

go run cmd/ingestor/main.go
go run cmd/nuevotwitter/main.go

Enviar mensajes:
curl -X POST localhost:8080/nombre%20del%20usuario/mi%20mensaje%20de%20twitter

Se pueden ver los mensajes almacenados en redis abriendo una sesion en el contenedor y utilizando redis-cli. Con keys * ver las claves de mensajes y con smembers ver el contenido.

Se pueden ver los mensajes que se van enconlando en kafka abriendo una sesion en el contenedor y corriendo el siguiente comando:
kafka-console-consumer --topic services --bootstrap-server localhost:9092 --property "print.key=true"
