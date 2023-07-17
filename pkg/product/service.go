package product

type Service interface {
	Find(map[string]string) ([]Product, error)
	Get(productID string) (*Product, error)
	Create(*Product) (*Product, error)
	Update(*Product) (*Product, error)
	Delete(string) (*Product, error)
	AddModifier(productID, modifierID string) (*Product, error)
	RemoveModifier(productID, modifierID string) (*Product, error)
	GetOverriders(productID, field string) ([]Overrider, error)

	ModifierFind(map[string]string) ([]Modifier, error)
	ModifierCreate(*Modifier) (*Modifier, error)
	ModifierAddProduct(productID, modifierID string) (*Modifier, error)
	ModifierRemoveProduct(productID, modifierID string) (*Modifier, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Find(filter map[string]string) ([]Product, error) {
	return s.repository.Find(filter)
}

func (s service) Get(productID string) (*Product, error) {
	return s.repository.Get(productID)
}

func (s service) Create(product *Product) (*Product, error) {
	return s.repository.Create(product)
}

func (s service) Update(product *Product) (*Product, error) {
	return s.repository.Update(product)
}

func (s service) Delete(productID string) (*Product, error) {
	return s.repository.Delete(productID)
}

func (s service) AddModifier(productID, modifierID string) (*Product, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.AddModifier(product, modifier)
}

func (s service) RemoveModifier(productID, modifierID string) (*Product, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.RemoveModifier(product, modifier)
}

func (s service) GetOverriders(productID, field string) ([]Overrider, error) {
	return s.repository.GetOverriders(productID, field)
}

func (s service) ModifierFind(filter map[string]string) ([]Modifier, error) {
	return s.repository.ModifierFind(filter)
}

func (s service) ModifierCreate(modifier *Modifier) (*Modifier, error) {
	return s.repository.ModifierCreate(modifier)
}

func (s service) ModifierAddProduct(productID, modifierID string) (*Modifier, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.ModifierAddProduct(product, modifier)
}

func (s service) ModifierRemoveProduct(productID, modifierID string) (*Modifier, error) {
	product, err := s.repository.Get(productID)
	if err != nil {
		return nil, err
	}

	modifier, err := s.repository.ModifierGet(modifierID)
	if err != nil {
		return nil, err
	}

	return s.repository.ModifierRemoveProduct(product, modifier)
}
