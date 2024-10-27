# TDL_Dinamita

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


## Ejecutar docker-compose
```
docker-compose up --build
```


### TODO

- Validar que el ingreso de la password no sea mayor a 72 bytes (ya que el hasheo con bycrypt solo acepta strings menores o iguales a 72 bytes).
    > GenerateFromPassword does not accept passwords longer than 72 bytes.