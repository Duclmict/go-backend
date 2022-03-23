package mail_service

import (
	"fmt"
    gomail "gopkg.in/gomail.v2"

	"github.com/Duclmict/go-backend/config"
	"github.com/Duclmict/go-backend/app/helper"
	"github.com/Duclmict/go-backend/app/service/log_service"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ses"
    "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	SES_DRIVER  = "ses"
	SMTP_DRIVER = "smtp"
)

var (
	CharSet = "UTF-8"
)

// public
func Send(template *config.MailTemplate, mail_address []*string, m_body string, cc_mail []config.CCStruct, attach []config.AttachStruct) error {

	log_service.Debug("Template:" + fmt.Sprint(template.From + "  " + template.Name + "  " + template.Subject))
	log_service.Debug("Address:" + fmt.Sprint(mail_address))
	log_service.Debug("body:" + fmt.Sprint(m_body))

	if (template.From == "" ||  len(mail_address) <= 0 || m_body == "") {
		return helper.ErrMailSettingError
	}

    switch config.MAIL_DRIVER {
		case SMTP_DRIVER:
			m := gomail.NewMessage()
			m.SetHeader("From", template.From, template.Name)
			for _, element := range mail_address {
				if element != nil {
					m.SetHeader("To", *element)
				}
			}
			
			m.SetHeader("Subject", template.Subject)
			m.SetBody("text/html", m_body)

			if (len(cc_mail) > 0) {
				for _, element_cc := range cc_mail {
					m.SetAddressHeader("Cc", element_cc.Address, element_cc.Name)
				}
			}

			if (len(attach) > 0) {
				for _, element_a := range attach {
					m.Attach(element_a.FileName)
				}
			}

			d := gomail.NewDialer(config.MAIL_HOST, config.MAIL_PORT, config.MAIL_USERNAME, config.MAIL_PASSWORD)

			// Send the email to Bob, Cora and Dan.
			if err := d.DialAndSend(m); err != nil {
				log_service.Error("Send mail error:" + fmt.Sprint(err))
				return err
			}
			
		case SES_DRIVER:
			// Create a new session in the us-west-2 region.
			// Replace us-west-2 with the AWS Region you're using for Amazon SES.
			sess, err := session.NewSession(&aws.Config{
					Region:aws.String(config.AWS_DEFAULT_REGION),
					Credentials: credentials.NewStaticCredentials(config.AWS_ACCESS_KEY_ID, config.AWS_SECRET_ACCESS_KEY, "TOKEN"),
				})
			
			// Create an SES session.
			svc := ses.New(sess)
			
			// Assemble the email.
			input := &ses.SendEmailInput{
				Destination: &ses.Destination{
					CcAddresses: []*string{
					},
					ToAddresses: mail_address,
				},
				Message: &ses.Message{
					Body: &ses.Body{
						Html: &ses.Content{
							Charset: aws.String(CharSet),
							Data:    aws.String(m_body),
						},
					},
					Subject: &ses.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(template.Subject),
					},
				},
				Source: aws.String(template.From),
					// Uncomment to use a configuration set
					//ConfigurationSetName: aws.String(ConfigurationSet),
			}

			// Attempt to send the email.
			_, err = svc.SendEmail(input)
			
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
	}

	return nil
}