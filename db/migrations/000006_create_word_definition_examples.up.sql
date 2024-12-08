CREATE TABLE word_definition_examples (
    id INTEGER PRIMARY KEY,
    word_definition_id INTEGER NOT NULL,
    example TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (word_definition_id) REFERENCES word_definitions(id) ON DELETE CASCADE
);