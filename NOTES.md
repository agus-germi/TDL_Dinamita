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
POST > localhost:8080/tables/register
```
{
  "number": 9,
  "seats": 4,
  "location": "hall"
}
```