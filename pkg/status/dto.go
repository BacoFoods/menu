package status

type CreateStatus struct {
	Name   string `json:"name" binding:"required"`
	Color  string `json:"color"`
	Code   string `json:"code" binding:"required"`
	Active bool   `json:"active"`
}

func (cs CreateStatus) ToStatus() *Status {
	return &Status{
		Name:   cs.Name,
		Color:  cs.Color,
		Code:   cs.Code,
		Active: cs.Active,
	}
}

type UpdateStatus struct {
	ID    uint   `json:"id" binding:"required"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Code  string `json:"code"`
}

func (us UpdateStatus) ToStatus() *Status {
	return &Status{
		ID:    us.ID,
		Name:  us.Name,
		Color: us.Color,
		Code:  us.Code,
	}
}
