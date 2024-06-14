-- table des utilisateurs
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(30) NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Création de la table des catégories
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Création de la table des messages de forum
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    title VARCHAR(150) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Insertion d'un utilisateur exemple
INSERT INTO users (username, email, password) VALUES ('john_doe', 'john@example.com', 'password123');

-- Insertion d'une catégorie exemple
INSERT INTO categories (name, description) VALUES ('General Discussion', 'Talk about anything and everything.');

-- Insertion d'un post exemple
INSERT INTO posts (user_id, category_id, title, content) VALUES (1, 1, 'Welcome to the Forum!', 'This is the first post on this forum.');
