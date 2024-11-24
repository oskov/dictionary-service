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

type WordRepository struct {
	db *sqlx.DB
}

func NewWordRepository(db *sqlx.DB) *WordRepository {
	return &WordRepository{db: db}
}

type DefinitionWithPriority struct {
	Definition string
	Priority   int64
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

	insertQuery := "INSERT INTO word_definitions (word_id, definition, priority) VALUES "

	insertQueryParts := make([]string, 0, len(definitions))
	insertValues := make([]any, 0, len(definitions)*3)

	for _, def := range definitions {
		insertQueryParts = append(insertQueryParts, "(?, ?, ?)")
		insertValues = append(insertValues, result.WordID, def.Definition, def.Priority)
	}

	insertQuery += strings.Join(insertQueryParts, ", ")

	insertQuery += "RETURNING id;"

	var ids []int64

	err = r.db.Select(&ids, insertQuery, insertValues...)
	if err != nil {
		return result, err
	}
	result.DefinitionIDs = ids

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
