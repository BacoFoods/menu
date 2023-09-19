package client

type Service interface {
	List() ([]Client, error)
	Get(id string) (*Client, error)
	Create(client *Client) (*Client, error)
	Update(client *Client) (*Client, error)
	Delete(id string) (*Client, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) List() ([]Client, error) {
	return s.repository.List()
}

func (s service) Get(id string) (*Client, error) {
	return s.repository.Get(id)
}

func (s service) Create(client *Client) (*Client, error) {
	return s.repository.Create(client)
}

func (s service) Update(client *Client) (*Client, error) {
	return s.repository.Update(client)
}

func (s service) Delete(id string) (*Client, error) {
	return s.repository.Delete(id)
}
