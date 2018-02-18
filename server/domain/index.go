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
	TgID   string `json:"tgid" validate:"required"`
	Title  string `json:"title"`
	Info   string `json:"info"`
	Desc   string `json:"desc"`
	Type   string `json:"type" validate:"required"`
	Tags   []Tag  `json:"tags" validate:"dive"`
	Imgsrc string `json:"imgsrc"`
}

// TgTagBucket ...
type TgTagBucket struct {
	Key      string `json:"key"`
	DocCount int32  `json:"doc_count"`
}

// Buckets ...
type TgTagBuckets struct {
	Buckets []TgTagBucket `json:"buckets"`
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

// NewTgTagBuckets ...
func NewTgTagBuckets() *TgTagBuckets {
	return &TgTagBuckets{}
}
