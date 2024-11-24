package core

import (
	"github.com/jmoiron/sqlx"
	"github.com/oskov/cambridge-dictionary-parser/parser"
	"github.com/oskov/dictionary-service/internal/core/repositories"
	"github.com/oskov/dictionary-service/internal/core/services"
	"github.com/oskov/dictionary-service/internal/util/lock"
)

type Core struct {
	WordService *services.WordService
	UserService *services.UserService

	wordLocker *lock.LockStorage[string]
}

func NewCore(db *sqlx.DB) *Core {
	wordRepo := repositories.NewWordRepository(db)
	userRepo := repositories.NewUserRepository(db)

	parser := parser.NewDictionaryParser()
	wordLocker := lock.NewLockStorage[string]()

	wordService := services.NewWordService(wordRepo, parser, wordLocker)
	userService := services.NewUserService(userRepo)

	return &Core{
		WordService: wordService,
		UserService: userService,
	}
}

func (c *Core) GetWord(word string) (services.GetWordResult, error) {
	return c.WordService.GetWord(word)
}

func (c *Core) Close() error {
	c.wordLocker.Close()
	return nil
}
