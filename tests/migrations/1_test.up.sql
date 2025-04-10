CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(50) PRIMARY KEY NOT NULL,
    min_amount INT NOT NULL
);

INSERT INTO users (user_id, min_amount) VALUES ('stray228', 5000);