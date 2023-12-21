CREATE TYPE role AS ENUM ('Buyer', 'Seller');

CREATE TABLE users (
    ID UUID PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL,
    hashedPassword VARCHAR(255) NOT NULL,
    email VARCHAR(320) UNIQUE NOT NULL,
    role role NOT NULL
);

