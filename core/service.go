package core

import (
	"errors"

	"github.com/ckalagara/group-a-accounts/model"
)

type Service interface {
	UpdateAccount(a model.Account) (model.Account, error)
	GetAccount(id string) (model.Account, error)
	DeleteAccount(id string) error
}

type service struct {
	store Store
}

func NewAccountService(s Store) Service {
	return &service{store: s}
}

func (s *service) UpdateAccount(a model.Account) (model.Account, error) {
	if a.ID == "" {
		return model.Account{}, errors.New("account ID cannot be empty")
	}

	updated, err := s.store.Update(a)
	if err != nil {
		return model.Account{}, err
	}
	return updated, nil
}

func (s *service) GetAccount(id string) (model.Account, error) {
	if id == "" {
		return model.Account{}, errors.New("account ID is required")
	}
	return s.store.Get(id)
}

func (s *service) DeleteAccount(id string) error {
	if id == "" {
		return errors.New("account ID is required")
	}
	return s.store.Delete(id)
}
