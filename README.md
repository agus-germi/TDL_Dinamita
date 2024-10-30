# TDL_Dinamita

## Run the following command in to execute the application:
```
docker-compose up --build
```

### TODO

- Validar que el ingreso de la password no sea mayor a 72 bytes (ya que el hasheo con bycrypt solo acepta strings menores o iguales a 72 bytes).
    > GenerateFromPassword does not accept passwords longer than 72 bytes.

- Create UserSession inside `entity` module. This struct will replace `dtos.User struct`.

- Migrate from sqlx to pgxpool

- **Reestructuracion**: No seria mala idea tener un servicio y un repositorio para cada entidad que forma parte de nuestro modelo de negocios. --> Consultar con el profe.
    > De esta forma encapsulariamos responsabilidades y cambios de requirimientos frente a cada entidad en particular.
    La desventaja es que aumenta bastante la complejidad de nuestro proyecto.
    > Hasta ahora lo que plantie quizas que es trampa pq simplemente cree archivos diferentes para las funciones que estan
    dirijidas a cada entidad en particular pero dichas funciones las implementa la misma estructura (repo struct)


## Useful Go commands:
- ```go mod tidy```
    > Ejecutarlo regularmente para actualizar las dependencias y eliminar aquellas que son innecesarias.
- ```go clean -modcache```
    > Se usa en casos especiales para limpiar la cache de Go

## Built information
Primera linea a ejecutar al crear el proyecto en go:
```
go mod init github.com/agus-germi/TDL_Dinamita
```

Para levantar el servidor:
```
go run ./cmd/server
```

Crear la imagen del contenedor.
```
docker build --tag restaurant_system
```

To see the list of images:
```
docker images
```

To run the image inside of a container:
```
docker run restaurant_system
```

Para hacer mas liviana la imagen de docker --> Multi-stage Dockerfile builds.

When you run this command, you’ll notice that you weren't returned to the command prompt. This is because your application is a REST server and will run in a loop waiting for incoming requests without returning control back to the OS until you stop the container.

Run the restaurant_system image into a container in the background and publish the internal port 8080 to 8080 in the host.
```
docker run -d -p 8080:8080 restaurant_system
d75e61fcad1e0c0eca69a3f767be6ba28a66625ce4dc42201a8a323e8313c14e
```
To name a container, you must pass the `--name` flag to the `run` command:
```
docker run -d -p 8080:8080 --name rest_server restaurant_system
d75e61fcad1e0c0eca69a3f767be6ba28a66625ce4dc42201a8a323e8313c14e
```

Again, make sure that your container is running. Run the same curl command:
```
curl http://localhost:8080/
Hello, Docker! <3
```

to list running containers:
```
docker ps
```
When you ran the docker ps command, the default output is to only show running containers. If you pass the --all or -a for short, you will see all containers on your system, including stopped containers and running containers.

to stop a container:
```
docker stop <container name or ID>
```

to restart running a container:
```
docker restart <container name or ID>
```

to remove a container:
```
docker rm <container_name>
```

Para instalarse Mockery y que se cree de forma automatica los mocks de las interfaces:
```
go install github.com/vektra/mockery/v2@latest
```

Para generar los mocks:
1. comentar arriba de la interfaz que se quiere mockear con: 
    >//go:generate mockery --name=<Interface_name> --output=<interface_name> --inpackage
    
    Notar que el primer nombre `Interface_name` comienza en mayuscula, mientras que el segundo `interface_name` en minúscula. (Esta claro que el 'interface_name' hay que cambiarlo por el nombre de nuestra interface)

2. Ejecutar el comando:
    ```
    go generate ./...
    ```

## To send request and tests if the endpoints are working execute:

- ```curl.exe -v localhost:8080```
    > This command is useful for Windows OS. If your host OS is Linux or MacOS run `curl` instead of `curl.exe`.
    Note: `-v` is a tag use to print a bunch of info about the request.






## Theoretical concepts
- [Pool de conexiones](https://chatgpt.com/share/671f899c-bc00-8004-acf7-133575a8e903)
    > Para configurar de forma mas precisa el pool de conexiones podemos migrar de `sqlx` a `pgxpool` 

- [REDIS](https://chatgpt.com/share/671fcb41-b638-8004-8635-137baf717321)

- [Middlewares](https://chatgpt.com/share/671fcfb4-23d4-8004-b223-6c34afdc5cf0)

- [Frameworks HTTP en Go](https://chatgpt.com/share/671fd0ea-363c-8004-9700-e3298ce46e8c)

- [Goose]() --> INVESTIGAR: go get github.com/pressly/goose/v3


## To keep in mind
- Al momento de desplegar la app en produccion deberiamos considerar usar herramientas de gestion de secretos (Docker secrets).
    > Para entornos de producción, considera usar Docker Secrets o herramientas de gestión de secretos, como AWS Secrets Manager o HashiCorp Vault, para manejar información sensible de manera más segura.
