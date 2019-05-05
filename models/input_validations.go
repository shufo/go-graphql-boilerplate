package models

import (
	"context"

	"github.com/shufo/go-graphql-boilerplate/translations"

	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/is"
)

func (i CreateUserInput) Validate(ctx context.Context) validation.Errors {
	errors := validation.Errors{
		translations.T(ctx, "email"): validation.Validate(i.Email,
			validation.Required.Error(translations.T(ctx, "required")),
			is.Email.Error(translations.T(ctx,
				"email_validation",
			)),
		),
		translations.T(ctx, "password"): validation.Validate(i.Password,
			validation.Required.Error(translations.T(ctx, "required")),
			validation.Length(6, 1024).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 6, "Max": 1024}),
			)),
		translations.T(ctx, "first_name"): validation.Validate(i.FirstName,
			validation.Required.Error(translations.T(ctx, "required")),
			validation.Length(1, 255).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 1, "Max": 255}),
			)),
		translations.T(ctx, "last_name"): validation.Validate(i.LastName,
			validation.Required.Error(translations.T(ctx, "required")),
			validation.Length(1, 255).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 1, "Max": 255}),
			)),
		translations.T(ctx, "phone_number"): validation.Validate(i.PhoneNumber,
			validation.Required.Error(translations.T(ctx, "required")),
			validation.Length(1, 15).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 1, "Max": 15}),
			)),
	}

	if errors.Filter() != nil {
		return errors
	}

	return nil
}

func (i AuthUserInput) Validate(ctx context.Context) validation.Errors {
	errors := validation.Errors{
		translations.T(ctx, "email"): validation.Validate(i.Email,
			validation.Required,
			is.Email.Error(translations.T(ctx,
				"email_validation",
			)),
		),
		translations.T(ctx, "password"): validation.Validate(i.Password,
			validation.Required,
			validation.Length(6, 1024).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 6, "Max": 1024}),
			)),
	}

	if errors.Filter() != nil {
		return errors
	}

	return nil
}

func (i RequestPasswordResetInput) Validate(ctx context.Context) validation.Errors {
	errors := validation.Errors{
		translations.T(ctx, "email"): validation.Validate(i.Email,
			validation.Required.Error(translations.T(ctx, "required")),
			is.Email.Error(translations.T(ctx,
				"email_validation",
			)),
		),
	}

	if errors.Filter() != nil {
		return errors
	}

	return nil
}

func (i ValidatePasswordResetInput) Validate(ctx context.Context) validation.Errors {
	errors := validation.Errors{
		translations.T(ctx, "token"): validation.Validate(i.Token,
			validation.Required.Error(translations.T(ctx, "required")),
			validation.Length(10, 100).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 10, "Max": 100}),
			)),
	}

	if errors.Filter() != nil {
		return errors
	}

	return nil
}

func (i CompletePasswordResetInput) Validate(ctx context.Context) validation.Errors {
	errors := validation.Errors{
		translations.T(ctx, "token"): validation.Validate(i.Token,
			validation.Required.Error(translations.T(ctx, "required")),
			validation.Length(10, 100).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 10, "Max": 100}),
			)),
		translations.T(ctx, "new_password"): validation.Validate(i.NewPassword,
			validation.Required.Error(translations.T(ctx, "required")),
			validation.Length(6, 1024).Error(translations.TWithTemplateData(ctx,
				"length_validation",
				map[string]interface{}{"Min": 6, "Max": 1024}),
			)),
	}

	if errors.Filter() != nil {
		return errors
	}

	return nil
}
