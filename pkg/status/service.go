package status

import "fmt"

type Service interface {
	Create(status *Status) (*Status, error)
	Delete(statusID string) error
	Get(statusID string) (*Status, error)
	Find() ([]Status, error)
	Update(*Status) (*Status, error)
	GetFirst() (*Status, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(status *Status) (*Status, error) {
	status_, err := s.repository.Create(status)
	if err != nil {
		return nil, err
	}

	if status_.PrevStatusID != nil {
		prevID := fmt.Sprintf("%v", *status_.PrevStatusID)

		prev, err := s.repository.Get(prevID)
		if err != nil {
			return nil, err
		}

		prev.Next = status_
		prev.NextStatusID = &status_.ID
		if _, err := s.repository.Update(prev, prevID); err != nil {
			return nil, err
		}

		status_.Prev = prev
	}

	if status_.NextStatusID != nil {
		nextID := fmt.Sprintf("%v", *status_.NextStatusID)

		next, err := s.repository.Get(nextID)
		if err != nil {
			return nil, err
		}

		next.Prev = status_
		next.PrevStatusID = &status_.ID
		if _, err := s.repository.Update(next, nextID); err != nil {
			return nil, err
		}

		status_.Next = next
	}

	return s.repository.Update(status_, fmt.Sprint(status_.ID))
}

func (s service) Delete(statusID string) error {
	// TODO: improve delete process to update prev and next status, define what happens when chain breaks
	return s.repository.Delete(statusID)
}

func (s service) Get(statusID string) (*Status, error) {
	return s.repository.Get(statusID)
}

func (s service) Find() ([]Status, error) {
	return s.repository.Find()
}

func (s service) GetFirst() (*Status, error) {
	return s.repository.GetFirst()
}

func (s service) Update(status *Status) (*Status, error) {
	return s.repository.Update(status, fmt.Sprint(status.ID))
}
