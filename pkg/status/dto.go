package status

type CreateStatus struct {
	Name         string `json:"name" binding:"required"`
	Color        string `json:"color"`
	Code         string `json:"code"`
	Active       bool   `json:"active"`
	PrevStatusID *uint  `json:"prev_status_id"`
	NextStatusID *uint  `json:"next_status_id"`
}

func (cs CreateStatus) ToStatus() *Status {
	return &Status{
		Name:         cs.Name,
		Color:        cs.Color,
		Code:         cs.Code,
		Active:       cs.Active,
		PrevStatusID: cs.PrevStatusID,
		NextStatusID: cs.NextStatusID,
	}
}

type UpdateStatus struct {
	ID           uint   `json:"id" binding:"required"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	Code         string `json:"code"`
	PrevStatusID *uint  `json:"prev_status_id"`
	NextStatusID *uint  `json:"next_status_id"`
}

func (us UpdateStatus) ToStatus() *Status {
	return &Status{
		ID:           us.ID,
		Name:         us.Name,
		Color:        us.Color,
		Code:         us.Code,
		PrevStatusID: us.PrevStatusID,
		NextStatusID: us.NextStatusID,
	}
}
