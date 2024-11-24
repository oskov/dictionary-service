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
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word VARCHAR(255) NOT NULL UNIQUE
    );
    CREATE TABLE word_definitions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word_id INTEGER NOT NULL,
        definition TEXT NOT NULL,
        priority INTEGER NOT NULL,
        FOREIGN KEY (word_id) REFERENCES words(id)
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
		{Definition: "definition1", Priority: 1},
		{Definition: "definition2", Priority: 2},
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
}
