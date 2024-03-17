CREATE TABLE IF NOT EXISTS actor (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gender VARCHAR(6) CHECK (gender IN ('male', 'female', 'other')) NOT NULL,
    birth_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS movie (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 150) NOT NULL,
    description TEXT CHECK (LENGTH(description) <= 1000),
    release_date DATE CHECK (release_date >= '1900-01-01') NOT NULL,
    rating INT CHECK (rating >= 0 AND rating <= 10) NOT NULL
);

CREATE TABLE IF NOT EXISTS movie_actor (
    movie_id INT,
    actor_id INT,
    FOREIGN KEY (movie_id) REFERENCES movie(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actor(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(5) CHECK (role IN ('admin', 'user')) NOT NULL
);

INSERT INTO users (username, password_hash, role) VALUES ('admin', '646a6664666e6e766e66626e765116c28e651a19013822c09e5c70c9fc425a66dc', 'admin');
