package repositories

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type WordDefinitionToUser struct {
	ID               int64 `db:"id"`
	WordDefinitionID int64 `db:"word_definition_id"`
	UserID           int64 `db:"user_id"`
}

type User struct {
	ID int64 `db:"id"`
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser() (int64, error) {
	query := "INSERT INTO users (id) VALUES (DEFAULT)"

	res, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetUserByID(id int64) (*User, error) {
	var user User
	query := "SELECT id FROM users WHERE id = ?"
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) AddWordDefinitionsToUser(userID int64, wordDefinitionIDs []int64) error {
	insertQuery := "INSERT INTO word_definitions_to_users (user_id, word_definition_id) VALUES "
	insertQueryParts := make([]string, 0, len(wordDefinitionIDs))
	insertQueryValues := make([]any, 0, len(wordDefinitionIDs))

	for _, wordDefinitionID := range wordDefinitionIDs {
		insertQueryParts = append(insertQueryParts, "(?, ?)")
		insertQueryValues = append(insertQueryValues, userID, wordDefinitionID)
	}

	insertQuery += strings.Join(insertQueryParts, ", ")

	_, err := r.db.Exec(insertQuery, insertQueryValues...)
	return err
}
