package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oskov/cambridge-dictionary-parser/parser"
	"github.com/oskov/dictionary-service/internal/core/repositories"
	"github.com/oskov/dictionary-service/internal/util/lock"
)

type WordService struct {
	wordRepo   *repositories.WordRepository
	parser     *parser.DictionaryParser
	wordLocker *lock.LockStorage[string]
}

func NewWordService(
	wordRepo *repositories.WordRepository,
	parser *parser.DictionaryParser,
	wordLocker *lock.LockStorage[string],
) *WordService {
	return &WordService{
		parser:     parser,
		wordRepo:   wordRepo,
		wordLocker: wordLocker,
	}
}

type GetWordResult struct {
	Word        string
	WordID      int64
	Definitions []GetWordResultDefinition
}

type GetWordResultDefinition struct {
	WordDefinitionID int64
	Definition       string
}

func (s *WordService) GetWord(word string) (GetWordResult, error) {
	mu := s.wordLocker.GetMutex(word)
	locked := mu.TryLockWithTimeout(time.Second)
	if !locked {
		return GetWordResult{}, fmt.Errorf("word %s is being processed by another user", word)
	}
	defer mu.Unlock()

	def, err := s.wordRepo.GetDefinitionsForWord(word)
	if err != nil && err != sql.ErrNoRows {
		return GetWordResult{}, err
	}
	if len(def) == 0 || err == sql.ErrNoRows {
		return s.loadNewWord(word)

	}
	var result GetWordResult
	result.Word = word
	result.WordID = def[0].WordID

	for _, v := range def {
		result.Definitions = append(result.Definitions, GetWordResultDefinition{
			WordDefinitionID: v.ID,
			Definition:       v.Definition,
		})
	}

	return GetWordResult{}, nil
}

func (s *WordService) loadNewWord(word string) (GetWordResult, error) {
	definitions, err := s.parser.ParseWord(word)
	if err != nil {
		return GetWordResult{}, err
	}

	defs := make([]repositories.DefinitionWithPriority, 0, len(definitions.Definitions))

	for i, v := range definitions.Definitions {
		defs = append(defs, repositories.DefinitionWithPriority{
			Definition: v.Definition,
			Priority:   int64(i),
		})
	}

	addResult, err := s.wordRepo.AddWordWithDefinitions(word, defs)
	if err != nil {
		return GetWordResult{}, err
	}

	var result GetWordResult
	result.Word = word
	result.WordID = addResult.WordID

	for i, v := range defs {
		result.Definitions = append(result.Definitions, GetWordResultDefinition{
			WordDefinitionID: addResult.DefinitionIDs[i],
			Definition:       v.Definition,
		})
	}

	return result, nil
}
