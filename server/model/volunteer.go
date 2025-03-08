package model

type Volunteer struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Contact     string   `json:"contact"`
	Files       []string `json:"files"`
	Description string   `json:"description,omitempty"`
	CreateAt    int64    `json:"create_at,omitempty"`
	UpdateAt    int64    `json:"update_at,omitempty"`
	DeleteAt    int64    `json:"delete_at"`
}
