package repositories

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Word struct {
	ID        int64     `db:"id"`
	Word      string    `db:"word"`
	CreatedAt time.Time `db:"created_at"`
}

type WordDefinition struct {
	ID         int64     `db:"id"`
	WordID     int64     `db:"word_id"`
	Definition string    `db:"definition"`
	Priority   int64     `db:"priority"`
	CreatedAt  time.Time `db:"created_at"`
}

type WordDefinitionExample struct {
	ID               int64     `db:"id"`
	WordDefinitionID int64     `db:"word_definition_id"`
	Example          string    `db:"example"`
	CreatedAt        time.Time `db:"created_at"`
}

type WordRepository struct {
	db *sqlx.DB
}

func NewWordRepository(db *sqlx.DB) *WordRepository {
	return &WordRepository{db: db}
}

type DefinitionWithPriority struct {
	Definition string
	Priority   int64
	Examples   []string
}

type AddResult struct {
	WordID        int64
	DefinitionIDs []int64
}

func (r *WordRepository) AddWordWithDefinitions(
	word string,
	definitions []DefinitionWithPriority,
) (AddResult, error) {
	var result AddResult
	err := r.db.Get(&result.WordID, "SELECT id FROM words WHERE word = ?", word)
	if err != nil {
		if err == sql.ErrNoRows {
			res, err := r.db.Exec("INSERT INTO words (word) VALUES (?)", word)
			if err != nil {
				return result, err
			}
			wordID, err := res.LastInsertId()
			if err != nil {
				return result, err
			}
			result.WordID = wordID
		} else {
			return result, err
		}
	}

	insertQuery := "INSERT INTO word_definitions (word_id, definition, priority, created_at) VALUES "

	insertQueryParts := make([]string, 0, len(definitions))
	insertValues := make([]any, 0, len(definitions)*3)

	for _, def := range definitions {
		insertQueryParts = append(insertQueryParts, "(?, ?, ?, ?)")
		insertValues = append(insertValues, result.WordID, def.Definition, def.Priority, time.Now())
	}

	insertQuery += strings.Join(insertQueryParts, ", ")
	insertQuery += " RETURNING id;"

	var definitionIDs []int64
	err = r.db.Select(&definitionIDs, insertQuery, insertValues...)
	if err != nil {
		return result, err
	}
	result.DefinitionIDs = definitionIDs

	for i, defID := range definitionIDs {
		exampleInsertQuery := "INSERT INTO word_definition_examples (word_definition_id, example, created_at) VALUES "
		exampleInsertQueryParts := make([]string, 0, len(definitions[i].Examples))
		exampleInsertValues := make([]any, 0, len(definitions[i].Examples)*2)

		for _, example := range definitions[i].Examples {
			exampleInsertQueryParts = append(exampleInsertQueryParts, "(?, ?, ?)")
			exampleInsertValues = append(exampleInsertValues, defID, example, time.Now())
		}

		exampleInsertQuery += strings.Join(exampleInsertQueryParts, ", ")
		_, err := r.db.Exec(exampleInsertQuery, exampleInsertValues...)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (r *WordRepository) GetDefinitionsForWord(word string) ([]WordDefinition, error) {
	var wordID int64
	err := r.db.Get(&wordID, "SELECT id FROM words WHERE word = ?", word)
	if err != nil {
		return nil, err
	}

	var definitions []WordDefinition
	err = r.db.Select(&definitions, "SELECT id, word_id, definition, priority FROM word_definitions WHERE word_id = ? ORDER BY priority", wordID)
	if err != nil {
		return nil, err
	}

	return definitions, nil
}

func (r *WordRepository) GetWordByID(id int64) (*Word, error) {
	var word Word
	err := r.db.Get(&word, "SELECT * FROM words WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &word, nil
}
