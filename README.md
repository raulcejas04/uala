*Requerimientos

Versiones que utilicé
golang version go1.21.4
docker version 27.3.1
docker-compose v2.18.1

*Ejecución
*Levantar redis
Entrar en ./redis y correr "docker-composer up -d"
*Levantar kafka
Entrar en ./kafka y correr "docker-composer up -d"
PENDIENTE DE RESOLVER CON KAFKA la primera vez no va a levantar y va a crear subdirectorios debajo de ./data 
para que levante hay que darle permisos 777 con chmod -R ugo+rwx ./data y correr de nuevo docker-compose up -d
cuando se reinicia puede que haya que borrar el directorio ./data y de nuevo darle permisos
Entrar en el directorio twitter y levantar los 2 microservicios:
go run cmd/ingestor/main.go
go run cmd/feed/main.go

*Enviar mensajes:
curl -X POST localhost:8080/nombre%20del%20usuario/mi%20mensaje%20de%20twitter

*Test 
Ir al directorio twitter/tests y correr "go test -v"

*Ver contenido de Redis
Se pueden ver los mensajes almacenados en redis abriendo una sesion en el contenedor y utilizando redis-cli. 
Para ver las claves: "keys *" y los mensajes y con "smembers clave"

*Ver mensajes encolados en kafka mientras se encolan
Se pueden ver los mensajes que se van enconlando en kafka abriendo una sesion en el contenedor y corriendo el siguiente comando:
kafka-console-consumer --topic services --bootstrap-server localhost:9092 --property "print.key=true"
