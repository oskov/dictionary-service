package repositories

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupTestDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	schema := `
CREATE TABLE words (
    id INTEGER PRIMARY KEY,
    word VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE word_definitions (
    id INTEGER PRIMARY KEY,
    word_id INTEGER NOT NULL,
    definition TEXT NOT NULL,
    priority INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
);

CREATE TABLE word_definition_examples (
    id INTEGER PRIMARY KEY,
    word_definition_id INTEGER NOT NULL,
    example TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (word_definition_id) REFERENCES word_definitions(id) ON DELETE CASCADE
);
    `
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestAddWordWithDefinitions(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewWordRepository(db)

	word := "test"
	definitions := []DefinitionWithPriority{
		{Definition: "definition1", Priority: 1, Examples: []string{"example1", "example2"}},
		{Definition: "definition2", Priority: 2, Examples: []string{"example3", "example4"}},
	}

	result, err := repo.AddWordWithDefinitions(word, definitions)
	assert.NoError(t, err)

	var wordID int64
	err = db.Get(&wordID, "SELECT id FROM words WHERE word = $1", word)
	assert.NoError(t, err)
	assert.NotZero(t, wordID)
	assert.Equal(t, result.WordID, wordID)

	var defs []WordDefinition
	err = db.Select(&defs, "SELECT id, word_id, definition, priority FROM word_definitions WHERE word_id = $1 ORDER BY priority", wordID)
	assert.NoError(t, err)
	assert.Len(t, defs, 2)
	assert.Equal(t, definitions[0].Definition, defs[0].Definition)
	assert.Equal(t, definitions[0].Priority, defs[0].Priority)
	assert.Equal(t, result.DefinitionIDs[0], defs[0].ID)

	assert.Equal(t, definitions[1].Definition, defs[1].Definition)
	assert.Equal(t, definitions[1].Priority, defs[1].Priority)
	assert.Equal(t, result.DefinitionIDs[1], defs[1].ID)

	// test that examples are added
	var examples []string
	err = db.Select(&examples, "SELECT example FROM word_definition_examples WHERE word_definition_id = $1", defs[0].ID)
	assert.NoError(t, err)
	assert.Len(t, examples, 2)
	assert.Contains(t, examples, "example1")
	assert.Contains(t, examples, "example2")

	err = db.Select(&examples, "SELECT example FROM word_definition_examples WHERE word_definition_id = $1", defs[1].ID)
	assert.NoError(t, err)
	assert.Len(t, examples, 2)
	assert.Contains(t, examples, "example3")
	assert.Contains(t, examples, "example4")
}
