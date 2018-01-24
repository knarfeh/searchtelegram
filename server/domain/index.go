package domain

import (
	"gopkg.in/go-playground/validator.v9"
)

// Tag ...
type Tag struct {
	Count int32  `json:"count" validate:"required"`
	Name  string `json:"name" validate:"required"`
}

// TgResource ...
type TgResource struct {
	Name string `json:"name" validate:"required"`
	Info string `json:"info"`
	Desc string `json:"desc"`
	Type string `json:"type" validate:"required"`
	Tags []Tag  `json:"tags" validate:"dive"`
}

type (
	// CustomValidator ...
	CustomValidator struct {
		validator *validator.Validate
	}
)

// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewValidator ...
func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// NewTgResource ...
func NewTgResource() *TgResource {
	return &TgResource{}
}
