# TODO:
- Sacar logica repetida de los handlers
- Por algun motivo la asignacion de id para una reserva esta vinculada con la asignacion de id para mesas:
    > O sea, si creo una mesa la base de datos le asigna un id (digamos 1), y a continuacion, cuando hago una reserva nueva
      la base de datos le asigna, a la reserva, un id siguiente al que esta en la tabla "tables" (en este ejemplo: id=2)

- Implement transactions in repository functions.
  Luego de esto podemos:
  > Simular múltiples solicitudes concurrentes para asegurarte de que el sistema maneja adecuadamente los bloqueos y evita inconsistencias.
  wrk -t4 -c100 -d10s http://localhost:8080/api/reservations

  > Utiliza herramientas como ab (Apache Benchmark), wrk o k6 para simular múltiples solicitudes concurrentes y verificar el rendimiento de tu API.

- Implementar limite de tiempo para las queries a la base de datos (utilizando context)

- Cambiar las URIs de los endpoints para que sean mas RESTful.

# Testing Endpoints with Postman
## Endpoint reserve_table
```
{
  "user_id": 9,
  "name": "John Doe",
  "password": "securePass123",
  "email": "john.doe@example.com",
  "table_number": 1,
  "reservation_date": "2024-11-15T19:00:00Z"
}
```
## Endpoint remove_reservation
# TODO :
 - Eliminar mediante un "reservation ID"

Por ahora funciona mediante: 
DELETE > localhost:8080/reservations/remove
```
{
  "user_id": 9
}
```
## Endpoint add_table
Simpre que se agrega una mesa se deja como "available"
POST > localhost:8080/tables/register
```
{
  "number": 9,
  "seats": 4,
  "location": "eg de location"
}
```
## Endpoint remove_user
 ### TODO:
 - ver que se deberia de eliminar de ambas tablas > Users y User Role
DELETE > localhost:8080/users/remove
```
{
  "user_id": 9
}
```
## Endpoint remove_table
Simpre que se agrega una mesa se deja como "available"
DELETE > localhost:8080/tables/remove
Se elimina una mesa solo mediante el numero de mesa, porque asumimos que este numero es unico para cada mesa. No puede estar repetido
```
{
  "number": 9,
}
```