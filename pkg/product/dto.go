package product

type OverriderDTO struct {
	ID         string `json:"id"`
	ProductID  uint   `json:"product_id"`
	PlaceName  string `json:"place_name"`
	PlaceID    uint   `json:"place_id"`
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
	ApplyPrice  float64  `json:"apply_price" gorm:"precision:18;scale:2"`
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
