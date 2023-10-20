package course

type Service interface {
	Find(map[string]any) ([]Course, error)
	Get(string) (Course, error)
	Create(Course) (Course, error)
	Delete(string) (Course, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(filter map[string]any) ([]Course, error) {
	return s.repository.Find(filter)
}

func (s service) Get(id string) (Course, error) {
	return s.repository.Get(id)
}

func (s service) Create(course Course) (Course, error) {
	return s.repository.Create(course)
}

func (s service) Delete(id string) (Course, error) {
	return s.repository.Delete(id)
}
