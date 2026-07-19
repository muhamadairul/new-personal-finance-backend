package requests

// CategoryRequest holds category payload validation
type CategoryRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Icon  int    `json:"icon" validate:"required"`
	Color int64  `json:"color" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=income expense"`
}
