# TDL_Dinamita

## Built information

- Run the following command to execute the application (including PostgreSQL DB):
    ```sh
    docker-compose up --build
    ```
- Add `-d` to run it on the background:
     ```sh
    docker-compose up -d --build
    ```
- To shut down the server and the DB engine:
     ```sh
    docker-compose down
    ```
- To run ***only the rest server***
    1. Build a docker image from the Dockerfile
        ```sh
        docker build --tag restaurant-system:1.0.0 .
        ```
    2. Run the image inside the container:
        ```sh
        docker run restaurant-system:1.0.0
        ```
    > **Aclaracion**: Esto es por el momento. A medida que vamos completando versiones del server, nos conviene contruir una imagen para cada version e ir subiendolas a un repositorio en docker hub. De esta forma cualquier persona puede hacer `docker pull <namespace:tag>` y se descarga localmente la imagen, para luego correrla en un contenedor. De esa forma nos aseguramos de que este usando una version valida del server.

- To run the concurrency tests using k6 framework (make sure to run all the containers first):
    ```sh
    docker-compose exec k6 run /script.js
    ```

- To run the k6 container on interactive mode execute:
    ```sh
    docker-compose exec k6 sh
    ```

### TODO

- Validar que el ingreso de la password no sea mayor a 72 bytes (ya que el hasheo con bycrypt solo acepta strings menores o iguales a 72 bytes).
    > GenerateFromPassword does not accept passwords longer than 72 bytes.

- Create UserSession inside `entity` module. This struct will replace `models.User struct`.

- Migrate from sqlx to pgxpool

- **Reestructuracion**: No seria mala idea tener un servicio y un repositorio para cada entidad que forma parte de nuestro modelo de negocios. --> Consultar con el profe.
    > De esta forma encapsulariamos responsabilidades y cambios de requirimientos frente a cada entidad en particular.
    La desventaja es que aumenta bastante la complejidad de nuestro proyecto.
    > Hasta ahora lo que plantie quizas que es trampa pq simplemente cree archivos diferentes para las funciones que estan
    dirijidas a cada entidad en particular pero dichas funciones las implementa la misma estructura (repo struct)


## Useful *Go commands*:

- First command to execute in a Go project:
    ```sh
    go mod init github.com/agus-germi/TDL_Dinamita
    ```
    It creats the `go.mod` file that have all the dependencies which need our project.

- Update the dependecies, deleting those that are no longer needed by our project:
    ```sh
    go mod tidy
    ```
    It's a good practice to execute it often during the development phase.

- Clean the Go cache.
    ```sh
    go clean -modcache
    ```
    Only execute it in specific cases. Suggestion: Investigate when it's necessary.


## Useful *Docker commands*

- Build a docker image from a Dockerfile
    ```sh
    docker build --tag <tag-name:version> .
    ```
    El `.` indica que el `Dockerfile` se encuentra en el directorio actual donde se esta ejecuntando el comando desde la terminal. `--tag` puede ser reemplazado por `-t`.
    
- List all the docker images
    ```sh
    docker images
    ```

- Delete an especific image `IMAGE_ID`
    ```sh
    docker rmi IMAGE_ID
    ```
    If you need to force the deletion add `-f` to the command, as follows:
    ```sh
    docker rmi -f IMAGE_ID
    ```

- To delete all the `dangling docker images`.
    ```sh
    docker image prune
    ```
    A *dangling image* has neither `name` nor `tag`. It has both as `<none>`.
    
    `-f` tag is added to force the deletion without any need for confirmation:
    ```sh
    docker image prune -f
    ```

- Delete all the images:
    ```sh
    docker rmi $(docker images -q)
    ```
    If any image if already used by a container, add `-f`:
    ```sh
    docker rmi -f $(docker images -q)
    ```

- To run an especific image inside of a container:
    ```sh
    docker run <image_tag>
    ```
    You can run the last image executing:
    ```sh
    docker run
    ```



## Miscellaneous information

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
