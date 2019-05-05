package mail

import (
	"fmt"
	"os"

	//go get -u github.com/aws/aws-sdk-go

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "info@example.co.jp"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

type Mail struct {
	Recipient string
	Subject   string
	HTMLBody  string
	TextBody  string
}

func New(to string) *Mail {
	return &Mail{Recipient: to}
}

func (m *Mail) SetSubject(sub string) {
	m.Subject = sub
}

func (m *Mail) SetHTMLBody(body string) {
	m.HTMLBody = body
}

func (m *Mail) SetTextBody(body string) {
	m.TextBody = body
}

func (m *Mail) Validate() error {
	if m.HTMLBody == "" {
		return fmt.Errorf("require mail html body")
	}

	if m.TextBody == "" {
		return fmt.Errorf("require mail text body")
	}

	if m.Recipient == "" {
		return fmt.Errorf("require mail recipient")
	}

	if m.Subject == "" {
		return fmt.Errorf("require mail subject")
	}

	return nil
}

// Send sends mail with setting
func (m *Mail) Send() error {

	// validate
	if err := m.Validate(); err != nil {
		return err
	}

	// return if APP_ENV is local
	if env, found := os.LookupEnv("APP_ENV"); found && env == "local" {
		return nil
	}
	// func main() {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create an SES session.
	svc := ses.New(sess)
	// Recipient := email
	Recipient := m.Recipient

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(m.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(m.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(m.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(m.Subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return err
	}

	fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(result)

	return nil
}
