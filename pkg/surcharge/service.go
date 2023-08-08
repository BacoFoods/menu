package surcharge

type Service interface {
	Find(query map[string]string) ([]Surcharge, error)
	Get(id string) (*Surcharge, error)
	Create(surcharge *Surcharge) (*Surcharge, error)
	Update(id string, surcharge *Surcharge) (*Surcharge, error)
	Delete(id string) (*Surcharge, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

// Find to find surcharges by query
func (s service) Find(query map[string]string) ([]Surcharge, error) {
	return s.repository.Find(query)
}

// Get to get surcharge by id
func (s service) Get(id string) (*Surcharge, error) {
	return s.repository.Get(id)
}

// Create to create surcharge
func (s service) Create(surcharge *Surcharge) (*Surcharge, error) {
	return s.repository.Create(surcharge)
}

// Update to update surcharge
func (s service) Update(id string, surcharge *Surcharge) (*Surcharge, error) {
	return s.repository.Update(id, surcharge)
}

// Delete to delete surcharge
func (s service) Delete(id string) (*Surcharge, error) {
	return s.repository.Delete(id)
}
