package services

import "github.com/oskov/dictionary-service/internal/core/repositories"

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser() (int64, error) {
	id, err := s.userRepo.CreateUser()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UserService) GetUserByID(id int64) (*repositories.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) AddWordDefinitionsToUser(userID int64, wordDefinitionIDs []int64) error {
	err := s.userRepo.AddWordDefinitionsToUser(userID, wordDefinitionIDs)
	if err != nil {
		return err
	}
	return nil
}
