CREATE TABLE word_definitions_to_users (
    id INT AUTOINCREMENT PRIMARY KEY NOT NULL,
    word_definition_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

    UNIQUE(word_definition_id, user_id),

    FOREIGN KEY (word_definition_id) REFERENCES word_definitions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);