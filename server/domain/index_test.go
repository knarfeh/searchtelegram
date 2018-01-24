package domain_test

import (
	"github.com/knarfeh/searchtelegram/server/domain"
	. "gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestIndexValidation(t *testing.T) {
	validator := domain.NewValidator()

	testTable := []struct {
		TestName   string
		Input      *domain.TgResource
		ErrMessage string
	}{
		{
			"No name",
			&domain.TgResource{
				Type: "group",
				Tags: []domain.Tag{
					domain.Tag{
						Name:  "chinese",
						Count: 1,
					},
				},
			},
			"Key: 'Input.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			"No tag name",
			&domain.TgResource{
				Type: "group",
				Name: "gotname",
				Tags: []domain.Tag{
					domain.Tag{
						Count: 1,
					},
				},
			},
			"Key: 'Input.Tags[0].Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			"Pass without tag",
			&domain.TgResource{
				Type: "group",
				Name: "telegram",
			},
			"",
		},
	}

	for _, testItem := range testTable {
		err := validator.Validate(testItem)
		if err != nil {
			Equal(t, err.Error(), testItem.ErrMessage)
		}
	}
}
