package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shufo/go-graphql-boilerplate/mail"
	"github.com/shufo/go-graphql-boilerplate/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/volatiletech/null"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shufo/go-graphql-boilerplate/models"
	"github.com/shufo/go-graphql-boilerplate/translations"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (r *mutationResolver) RequestPasswordReset(ctx context.Context, input models.RequestPasswordResetInput) (*models.PasswordReset, error) {
	if err := input.Validate(ctx); err != nil {
		for k, v := range err {
			graphql.AddErrorf(ctx, "%s: %s", k, v)
		}
		return nil, nil
	}

	// get db instance
	db := ctx.Value("db").(*sql.DB)

	// check if user exists
	ap, err := models.AuthenticationProviders(
		qm.Where("provider_type = ?", "email"),
		qm.Where("provider_username = ?", input.Email),
	).One(ctx, db)

	if err != nil {
		return nil, fmt.Errorf(translations.T(ctx, "email_not_found"))
	}

	resetToken := utils.RandomUUID()

	pr := &models.PasswordReset{
		Status:             null.StringFrom("sent"),
		PasswordResetToken: null.StringFrom(resetToken),
		ExpiresAt:          null.TimeFrom(time.Now().Add(24 * time.Hour)),
	}

	if err := pr.SetAuthenticationProvider(ctx, db, false, ap); err != nil {
		return nil, err
	}

	if err := pr.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	// send reset email
	m := mail.New(input.Email)
	m.SetSubject(translations.T(ctx, "subject_password_reset"))

	variables := map[string]interface{}{
		"ResetLink": resetToken,
	}

	m.SetHTMLBody(translations.TWithTemplateData(ctx, "email_password_reset", variables))
	m.SetTextBody(translations.TWithTemplateData(ctx, "email_password_reset", variables))

	/* uncomment this if you want to really send email

	if err := m.Send(); err != nil {
		return nil, err
	}

	*/

	return pr, nil
}

func (r *mutationResolver) ValidatePasswordReset(ctx context.Context, input models.ValidatePasswordResetInput) (*models.PasswordReset, error) {
	// validate input
	if err := input.Validate(ctx); err != nil {
		for k, v := range err {
			graphql.AddErrorf(ctx, "%s: %s", k, v)
		}
		return nil, nil
	}

	// get db instance
	db := ctx.Value("db").(*sql.DB)

	// check if reset token exists
	pr, err := models.PasswordResets(
		qm.Where("password_reset_token = ?", input.Token),
		qm.Where("status = ?", "sent"),
		qm.Where("expires_at > ?", time.Now()),
	).One(ctx, db)

	if err != nil {
		return nil, fmt.Errorf(translations.T(ctx, "invalid_password_reset_token"))
	}

	pr.Status = null.StringFrom("verified")

	if _, err := pr.Update(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	return pr, nil
}

func (r *mutationResolver) CompletePasswordReset(ctx context.Context, input models.CompletePasswordResetInput) (*models.AuthenticationProvider, error) {
	// validate input
	if err := input.Validate(ctx); err != nil {
		for k, v := range err {
			graphql.AddErrorf(ctx, "%s: %s", k, v)
		}
		return nil, nil
	}

	// get db instance
	db := ctx.Value("db").(*sql.DB)

	// check if reset token exists
	pr, err := models.PasswordResets(
		qm.Where("password_reset_token = ?", input.Token),
		qm.Where("status = ?", "verified"),
		qm.Load("AuthenticationProvider"),
	).One(ctx, db)

	if err != nil {
		return nil, fmt.Errorf(translations.T(ctx, "verified_token_not_found"))
	}

	ap := pr.R.AuthenticationProvider

	// update authentication provider with new password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
	ap.ProviderPassword = string(hashed)

	if _, err := ap.Update(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	// send complete email
	m := mail.New(ap.Email.String)
	m.SetSubject(translations.T(ctx, "subject_password_reset_complete"))

	variables := map[string]interface{}{
		"email": ap.Email.String,
	}

	m.SetHTMLBody(translations.TWithTemplateData(ctx, "email_password_reset_complete", variables))
	m.SetTextBody(translations.TWithTemplateData(ctx, "email_password_reset_complete", variables))

	/* uncomment this if you want to really send email

	if err := m.Send(); err != nil {
		return nil, err
	}

	*/

	return ap, nil
}
