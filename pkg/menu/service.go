package menu

type Service interface {
	Create(*Menu) (*Menu, error)
	Find() ([]Menu, error)
	Get(string) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(menu *Menu) (*Menu, error) {
	return s.repo.Create(menu)
}

func (s *service) Find() ([]Menu, error) {
	return s.repo.Find()
}

func (s *service) Get(id string) (*Menu, error) {
	return s.repo.Get(id)
}

func (s *service) Update(menu *Menu) (*Menu, error) {
	return s.repo.Update(menu)
}

func (s *service) Delete(id string) (*Menu, error) {
	return s.repo.Delete(id)
}
