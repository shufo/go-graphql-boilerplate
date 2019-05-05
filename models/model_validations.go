package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Length(1, 255)),
	)
}
