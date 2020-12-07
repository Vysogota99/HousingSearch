CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    vk_profile VARCHAR(255) UNIQUE,
    telephone_number VARCHAR(15) UNIQUE NOT NULL,
    role INTEGER NOT NULL REFERENCES roles,
    password VARCHAR(63) NOT NULL,
    avatar_path text NOT NULL
);