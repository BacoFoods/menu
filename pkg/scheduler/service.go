package scheduler

type Service interface {
	Find(filter map[string]any) ([]Schedule, error)
	Create(schedule *Schedule) error
	Update(schedule *Schedule) error
	Delete(schedule *Schedule) error
	Today(storeID string) (*Schedule, error)
	TodayStores(brandID string) ([]Schedule, error)
	EnableStore(storeID string, enable bool) ([]Schedule, error)
	CreateHoliday(holiday *Holiday) (*Holiday, error)
	UpdateHoliday(holiday *Holiday) (*Holiday, error)
	DeleteHoliday(holiday *Holiday) error
	FindHoliday() ([]Holiday, error)
	GetTodayHoliday() (*Holiday, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(filter map[string]any) ([]Schedule, error) {
	return s.repository.Find(filter)
}

func (s service) Create(schedule *Schedule) error {
	return s.repository.Create(schedule)
}

func (s service) Update(schedule *Schedule) error {
	return s.repository.Update(schedule)
}

func (s service) Delete(schedule *Schedule) error {
	return s.repository.Delete(schedule)
}

func (s service) Today(storeID string) (*Schedule, error) {
	return s.repository.TodayStore(storeID)
}

func (s service) TodayStores(brandID string) ([]Schedule, error) {
	return s.repository.TodayBrand(brandID)
}

func (s service) EnableStore(storeID string, enable bool) ([]Schedule, error) {
	return s.repository.EnableStore(storeID, enable)
}

func (s service) CreateHoliday(holiday *Holiday) (*Holiday, error) {
	return s.repository.CreateHoliday(holiday)
}

func (s service) UpdateHoliday(holiday *Holiday) (*Holiday, error) {
	return s.repository.UpdateHoliday(holiday)
}

func (s service) DeleteHoliday(holiday *Holiday) error {
	return s.repository.DeleteHoliday(holiday)
}

func (s service) FindHoliday() ([]Holiday, error) {
	return s.repository.FindHoliday()
}

func (s service) GetTodayHoliday() (*Holiday, error) {
	return s.repository.GetTodayHoliday()
}
