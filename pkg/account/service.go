package account

type Service interface {
	Create(*Account) (*Account, error)
	Login(username, password string) (*Account, error)
	Delete(id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(account *Account) (*Account, error) {

	return s.repository.Create(account)
}

func (s service) Login(username, password string) (*Account, error) {

	return s.repository.Login(username, password)
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}
