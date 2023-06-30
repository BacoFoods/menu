package category

type Service interface {
	Find(map[string]string) ([]Category, error)
	Get(string) (*Category, error)
	Create(*Category) (*Category, error)
	Update(*Category) (*Category, error)
	Delete(string) (*Category, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(filter map[string]string) ([]Category, error) {
	return s.repository.Find(filter)
}

func (s service) Get(categoryID string) (*Category, error) {
	return s.repository.Get(categoryID)
}

func (s service) Create(category *Category) (*Category, error) {
	return s.repository.Create(category)
}

func (s service) Update(category *Category) (*Category, error) {
	return s.repository.Update(category)
}

func (s service) Delete(categoryID string) (*Category, error) {
	return s.repository.Delete(categoryID)
}
