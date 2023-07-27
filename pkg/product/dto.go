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
