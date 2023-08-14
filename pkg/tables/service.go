package tables

type Service interface {
	Get(id string) (*Table, error)
	Find(query map[string]any) ([]Table, error)
	Create(table *Table) (*Table, error)
	Update(id string, table *Table) (*Table, error)
	Delete(id string) error
	Release(id *uint) (*Table, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Get(id string) (*Table, error) {
	return s.repository.Get(id)
}

func (s service) Find(query map[string]any) ([]Table, error) {
	return s.repository.Find(query)
}

func (s service) Create(table *Table) (*Table, error) {
	return s.repository.Create(table)
}

func (s service) Update(id string, table *Table) (*Table, error) {
	return s.repository.Update(id, table)
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s service) Release(id *uint) (*Table, error) {
	return s.repository.RemoveOrder(id)
}
