CREATE TABLE spaced_repetition_cards (
    id INTEGER PRIMARY KEY,
    word_definition_to_user_id INTEGER NOT NULL,
    ease_factor DECIMAL(4, 2) DEFAULT 2.5,
    interval INTEGER NOT NULL,
    repetition INTEGER NOT NULL,
    due_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (word_definition_to_user_id) REFERENCES word_definitions_to_users(id) ON DELETE CASCADE
);