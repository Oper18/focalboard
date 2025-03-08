package model

type Role struct {
	// The role ID
	// required: true
	ID string `json:"id"`

	// The role name
	// required: true
	Name string `json:"name"`

	// Array of permissions for role
	Permissions []string `json:"permissions,omitempty"`
}
