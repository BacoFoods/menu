package menu

type Service interface {
	Find(map[string]string) ([]Menu, error)
	Get(string) (*Menu, error)
	Create(*Menu) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(filter map[string]string) ([]Menu, error) {
	return s.repository.Find(filter)
}

func (s service) Get(menuID string) (*Menu, error) {
	return s.repository.Get(menuID)
}

func (s service) Create(menu *Menu) (*Menu, error) {
	return s.repository.Create(menu)
}

func (s service) Update(menu *Menu) (*Menu, error) {
	return s.repository.Update(menu)
}

func (s service) Delete(menuID string) (*Menu, error) {
	return s.repository.Delete(menuID)
}
