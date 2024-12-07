--CREATE DATABASE tdl_dinamita;  Esta linea solo se deberia ejecutar una única vez.

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role_id INT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS tables (
    id SERIAL PRIMARY KEY,
    number INT NOT NULL UNIQUE,
    seats INT NOT NULL,
    location VARCHAR(255),    
    is_available BOOLEAN DEFAULT TRUE NOT NULL
);

CREATE TABLE IF NOT EXISTS time_slots (
    id SERIAL PRIMARY KEY,
    time TIME NOT NULL --  '12:00', '14:00', etc.
);

CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    reserved_by INT NOT NULL,
    table_number INT NOT NULL,
    date DATE NOT NULL, -- Reservation date
    time_slot_id INT NOT NULL,
    FOREIGN KEY (reserved_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE CASCADE,
    FOREIGN KEY (time_slot_id) REFERENCES time_slots(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS dishes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,-- para guardar precios de hasta 99999999.99
    description TEXT -- Descripcion del plato
);

CREATE TABLE IF NOT EXISTS opinions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT null,
    opinion TEXT,
    rating INT NOT null
);

-- Agregamos horarios fijos de manera dinamica 
-- Turnos desde las 12:00 hasta las 22:00 - considerando que cada turno es de 2hs
DO $$
BEGIN
    FOR hour IN 12..22 BY 2 LOOP
        INSERT INTO time_slots (time) VALUES (MAKE_TIME(hour, 0, 0));
    END LOOP;
END $$;


-- Creacion de los usuarios administradores. (En nuestro caso solo sera admin y user por ende no tendremos superadmin)
--  role 1 : Admin | role 2 : User

INSERT INTO roles (id, name) VALUES (1, 'Admin');
INSERT INTO roles (id, name) VALUES (2, 'User');

INSERT INTO users (name, password, email, role_id)
VALUES('agus', '$2a$14$qYw4gqHgSAObMdAraUZZ5.IVduI/6JJUze.WGYvs7fdViMyiOUBCG', 'agerminario@fi.uba.ar', 1);

INSERT INTO users (name, password, email, role_id) 
VALUES('valen', '$2a$14$fd5pMYtyaB2Z0pFogjwUjuHQY2I7PL23YTCM3.NI83KqncMATlwbS', 'vmoreno@fi.uba.ar', 1);

INSERT INTO users (name, password, email, role_id) 
VALUES('seba', '$2a$14$yeldkzqFM65K0qiCmlsCxuxBsEuKgwQThY89zz8M8NFURWLB1CH7u', 'skraglievich@fi.uba.ar', 1);