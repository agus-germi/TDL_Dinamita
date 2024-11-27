--CREATE DATABASE tdl_dinamita;  Esta linea solo se deberia ejecutar una única vez.

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL
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

-- Agregamos horarios fijos de manera dinamica 
-- Turnos desde  las 12:00 hasta las 22:00 - considerando que cada turno es de 2hs
DO $$
BEGIN
    FOR hour IN 12..22 BY 2 LOOP
        INSERT INTO time_slots (time) VALUES (MAKE_TIME(hour, 0, 0));
    END LOOP;
END $$;


-- Creacion de los usuarios administradores. (En nuestro caso solo sera admin y user por ende no tendremos superadmin)
--  role 1 : Admin | role 2 : User

INSERT INTO users (name, password, email, role) 
VALUES ('Agus', '12345678', 'agerminario@fi.uba.ar', 1);

INSERT INTO users (name, password, email, role) 
VALUES ('Valen', '12345678', 'vmorenofi.uba.ar', 1);

INSERT INTO users (name, password, email, role) 
VALUES ('Seba', '12345678', 'skraglievich@fi.uba.ar', 1);