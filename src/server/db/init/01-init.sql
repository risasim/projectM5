CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    isAdmin BOOLEAN DEFAULT FALSE,
    pi_SN VARCHAR(255),
    api_key VARCHAR(255)
    deathSound VARCHAR(255),
    );