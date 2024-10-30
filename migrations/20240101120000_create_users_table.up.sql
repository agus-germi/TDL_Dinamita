--CREATE DATABASE tdl_dinamita;  Esta linea solo se deberia ejecutar una única vez.

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tables (
    id SERIAL PRIMARY KEY,
    number INT NOT NULL UNIQUE,
    seats INT NOT NULL,
    location VARCHAR(255),    -- No seria mejor nombrarlo description?
    is_available BOOLEAN DEFAULT TRUE NOT NULL
);

CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    table_number INT NOT NULL,
    reserved_by INT NOT NULL,
    -- booking_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, --> Deseamos que haya un campo que indique el momento en que se realizo la reserva?
    reservation_date TIMESTAMP,   --Si quisieramos que tenga en cuenta la zona horaria tendriamos que usar el tipo de dato TIMESTAMPTZ
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE CASCADE,
    FOREIGN KEY (reserved_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_roles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,  -- Cada user solo puede estar vinculado a un único rol, por eso debe aparecer una única vez en la tabla user_roles
    role_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE  --Creo que tiene sentido que si se elimina un rol de la tabla roles,
                                                                -- entonces se debe eliminar todas las filas de la tabla users_roles que presenten el rol que se elimino.
);

INSERT INTO roles (name) VALUES ("admin")     -- Al ingresar los roles de esta forma: admin-->role_id=1 y customer-->role_id=2
INSERT INTO roles (name) VALUES ("customer")  -- or client
