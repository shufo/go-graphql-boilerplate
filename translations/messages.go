package translations

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

/***********
 * nouns
 ***********/

var name = i18n.Message{
	ID:          "name",
	Description: "The name of a person",
	One:         "Name",
	Other:       "Name",
}

var email = i18n.Message{
	ID:          "email",
	Description: "The email address of the user",
	One:         "Email",
	Other:       "Email",
}

var password = i18n.Message{
	ID:          "password",
	Description: "The passphrase of user",
	One:         "Password",
	Other:       "Password",
}

var fitst_name = i18n.Message{
	ID:          "first_name",
	Description: "The first name of user",
	One:         "First Name",
	Other:       "First Name",
}

var last_name = i18n.Message{
	ID:          "last_name",
	Description: "The last name of user",
	One:         "Last Name",
	Other:       "Last Name",
}

var phone_number = i18n.Message{
	ID:          "phone_number",
	Description: "The phone number of user",
	One:         "Phone Number",
	Other:       "Phone Number",
}

var password_verification = i18n.Message{
	ID:          "password_verification",
	Description: "The passphrase of user (verify)",
	One:         "Password",
	Other:       "Password",
}

var token = i18n.Message{
	ID:          "token",
	Description: "The token used by system",
	One:         "Token",
	Other:       "Tokens",
}

var new_password = i18n.Message{
	ID:          "new_password",
	Description: "A new password when user forget password",
	One:         "New password",
	Other:       "New password",
}

/***********
 * messages
 ***********/

var required = i18n.Message{
	ID:          "required",
	Description: "The message indicates input is required",
	One:         "cannot be blank",
	Other:       "cannot be blank",
}

var email_already_exists = i18n.Message{
	ID:          "email_already_exists",
	Description: "The message specified email already used",
	One:         "The specified email is already used",
	Other:       "The specified email is already used",
}

var email_not_found = i18n.Message{
	ID:          "email_not_found",
	Description: "The message specified email is not found",
	One:         "The specified email is not found",
	Other:       "The specified email is not found",
}

var user_not_found = i18n.Message{
	ID:          "user_not_found",
	Description: "The message specified user is not found",
	One:         "The specified user is not found",
	Other:       "The specified user is not found",
}

var email_or_password_is_incorrect = i18n.Message{
	ID:          "email_or_password_is_incorrect",
	Description: "The message for login failed",
	One:         "Email or Password is incorrect",
	Other:       "Email or Password is incorrect",
}

var length_validation = i18n.Message{
	ID:          "length_validation",
	Description: "The validation message of input length",
	One:         "Requires {{.Min}} to {{.Max}} characters",
	Other:       "Requires {{.Min}} to {{.Max}} characters",
}

var email_validation = i18n.Message{
	ID:          "email_validation",
	Description: "The validation message of email format",
	One:         "Requires email format",
	Other:       "Requires email format",
}

var field_match_validation = i18n.Message{
	ID:          "field_match_validation",
	Description: "The field must match with the other field",
	One:         "{{.One}} and {{.Other}} do not match",
	Other:       "{{.One}} and {{.Other}} do not match",
}

var invalid_password_reset_token = i18n.Message{
	ID:          "invalid_password_reset_token",
	Description: "The message when password reset token is invalid",
	One:         "Password reset token is invalid",
	Other:       "Password reset token is invalid",
}

var verified_token_not_found = i18n.Message{
	ID:          "verified_token_not_found",
	Description: "The message verified token is not found",
	One:         "There is no verified password reset token",
	Other:       "There is no verified password reset token",
}

/***********
 * Subjects
 ***********/

var subject_password_reset = i18n.Message{
	ID:          "subject_password_reset",
	Description: "The subject of password reset email",
	One:         "Change password for example",
	Other:       "Change password for example",
}

var subject_password_reset_complete = i18n.Message{
	ID:          "subject_password_reset_complete",
	Description: "The subject of password reset complete",
	One:         "Password reset complete",
	Other:       "Password reset complete",
}

/***********
 * Templates
 ***********/

var email_jpassword_reset = i18n.Message{
	ID:          "email_password_reset",
	Description: "The password reset email",
	One:         "<p>Sent by example.jp.</p>Please click {{.ResetLink}} in 24 hours.<br>Thank you.",
	Other:       "<p>Sent by example.jp.</p>Please click {{.ResetLink}} in 24 hours.<br>Thank you.",
}

var email_password_reset_complete = i18n.Message{
	ID:          "email_password_reset_complete",
	Description: "The password reset email",
	One:         "<p>Password reset complete.</p>.",
	Other:       "<p>Password reset complete.</p>",
}
