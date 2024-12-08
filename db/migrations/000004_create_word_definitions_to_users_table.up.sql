CREATE TABLE word_definitions_to_users (
    id INTEGER PRIMARY KEY,
    word_definition_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(word_definition_id, user_id),
    UNIQUE(word_definition_id, user_id),
    FOREIGN KEY (word_definition_id) REFERENCES word_definitions(id) ON DELETE CASCADE
);