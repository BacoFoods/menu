package status

import "fmt"

type Service interface {
	Create(status *Status) (*Status, error)
	Delete(statusID string) error
	Get(statusID string) (*Status, error)
	Find() ([]Status, error)
	Update(*Status) (*Status, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(status *Status) (*Status, error) {
	return s.repository.Create(status)
}

func (s service) Delete(statusID string) error {
	return s.repository.Delete(statusID)
}

func (s service) Get(statusID string) (*Status, error) {
	return s.repository.Get(statusID)
}

func (s service) Find() ([]Status, error) {
	return s.repository.Find()
}

func (s service) Update(status *Status) (*Status, error) {
	return s.repository.Update(status, fmt.Sprint(status.ID))
}
