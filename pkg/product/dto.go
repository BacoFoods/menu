package product

type Overrider struct {
	ID         string `json:"id"`
	PlaceName  string `json:"place_name"`
	FieldValue string `json:"field_value"`
}

type CategoryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ModifierDTO struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	ApplyPrice  float32  `json:"apply_price"`
	Category    Category `json:"category"`
}

func (dto ModifierDTO) ToModifier() Modifier {
	return Modifier{
		Name:        dto.Name,
		Description: dto.Description,
		Image:       dto.Image,
		ApplyPrice:  dto.ApplyPrice,
		Category:    dto.Category,
	}
}
