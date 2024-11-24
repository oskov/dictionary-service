CREATE TABLE spaced_repetition_cards (
    id INT AUTOINCREMENT PRIMARY KEY NOT NULL,
    word_definition_to_user_id INT NOT NULL,
    ease_factor FLOAT DEFAULT 2.5,
    interval INT NOT NULL,
    repetition INT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (word_definition_to_user_id) REFERENCES word_definitions_to_users(id) ON DELETE CASCADE 
);