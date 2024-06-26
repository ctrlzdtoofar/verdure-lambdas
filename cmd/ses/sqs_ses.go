package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/ctrlzdtoofar/verdure-lambdas/internal/mdl"
	"github.com/ctrlzdtoofar/verdure-lambdas/internal/settings"
	"strings"
)

// handleRequest processes SQSEvent messages for email confirmation actions,
// utilizing AWS SES for sending out emails. Each SQSEvent record is expected
// to contain a JSON message that details user confirmation data. This data
// is used to either confirm a new user's email address or reset a user's password.
//
// The function supports handling multiple records from an SQSEvent, deserializing
// each record's JSON into a UserConfirmation type, and sending an email using
// the appropriate SES template based on the confirmation type indicated in the message.
//
// Parameters:
//   - ctx: A context.Context to allow for timeout or cancellation signals to be
//     respected during execution.
//   - sqsEvent: An events.SQSEvent struct provided by the AWS SDK, containing one
//     or more records from an SQS message payload.
//
// Returns:
//   - An error if any issues occur during the processing of SQSEvent records, such
//     as failures in loading AWS configurations, deserializing JSON, or sending emails
//     via SES. Errors are wrapped with contextual information to aid in debugging.
func handleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load default config, %v", err)
	}

	sesSvc := sesv2.NewFromConfig(cfg)

	for _, message := range sqsEvent.Records {
		fmt.Printf("Recieved message %s from event source %s\n", message.Body, message.EventSource)

		confirm, err := mdl.UserConfirmationFromJson(message.Body)
		if err != nil {
			return fmt.Errorf("failed to deserialize confirmation json, %v", err)
		}

		fromAddress, err := settings.GetNoReplySecret(ctx, "noreply/Email")
		if err != nil {
			return fmt.Errorf("failed to get secret email from address, %v", err)
		}

		fmt.Printf("Using From Address %s\n", fromAddress)

		err = sendEmail(ctx, *sesSvc, *confirm, fromAddress)
		if err != nil {
			return err
		}
	}
	return nil
}

// sendEmail constructs and sends an email based on the user confirmation details using AWS SES.
// This function first determines the appropriate email template name by invoking the determineTemplate
// function, which formats the template name based on the confirmation type and language specified in the
// UserConfirmation struct. It then creates a template data string with a URL included from the UserConfirmation.
//
// The function sets up the email parameters, including the recipient's address, sender's address, and
// content, which consists of the chosen template and the template data. After setting up these parameters,
// it sends the email via the provided SES service.
//
// Parameters:
//
//	ctx - The context for controlling cancellations and timeouts.
//	sesSvc - The SES v2 client used to send the email.
//	confirm - A struct containing user confirmation details including the type of confirmation,
//	          the base URL for the confirmation link, user's email, etc.
//	from - The email address from which the email is sent.
//
// Returns:
//
//	err - An error value which is non-nil in case of failures in sending the email. If the email
//	      is successfully sent, the function returns nil.
//
// Usage:
//
//	This function is typically used in account creation or password reset flows where it is necessary
//	to send a verification or reset link to the user's email address.
func sendEmail(ctx context.Context, sesSvc sesv2.Client, confirm mdl.UserConfirmation, from string) (err error) {

	templateName := determineTemplate(confirm.ConfirmationType, confirm.Lang)
	templateData := fmt.Sprintf("{\"url\": \"%s\"}", confirm.ConfirmUrL())

	// Define Email Parameters
	input := &sesv2.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				confirm.Email,
			},
		},
		FromEmailAddress: &from,
		Content: &types.EmailContent{
			Template: &types.Template{
				TemplateData: &templateData,
				TemplateName: &templateName,
			},
		},
		EmailTags: []types.MessageTag{
			{
				Name:  aws.String("email_type"),
				Value: aws.String("confirmation"),
			},
		},
	}

	// Send Email
	_, err = sesSvc.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send %s email, %v", templateName, err)
	} else {
		fmt.Printf("Email successfully sent to %s\n", confirm.Email)
	}

	return
}

// determineTemplate generates the template name based on the confirmation type and language.
// This function supports two types of confirmation: NewUser and PasswordReset. It formats
// the template name by appending the properly cased language code to the base template name.
// For 'NewUser', it prefixes with "EmailConfirmation", and for password resets, it uses
// "PasswordResetConfirmation". If the language string is less than 2 characters long, it defaults
// to the English template without a language suffix.
//
// Parameters:
//
//	confirmType - The type of confirmation, which determines the base part of the template name.
//	lang - The ISO language code to append to the template name. The first letter will be
//	       uppercase and the rest lowercase if the length is greater than 1; otherwise, it uses
//	       the base template name.
//
// Returns:
//
//	templateName - The formatted template name as a string.
func determineTemplate(confirmType mdl.ConfirmationType, lang string) (templateName string) {
	if confirmType == mdl.NewUser {
		if len(lang) > 1 {
			templateName = fmt.Sprintf("EmailConfirmation%s",
				strings.ToUpper(lang[:1])+strings.ToLower(lang[1:]))
		} else {
			templateName = "EmailConfirmation"
		}
	} else {
		if len(lang) > 1 {
			templateName = fmt.Sprintf("PasswordResetConfirmation%s",
				strings.ToUpper(lang[:1])+strings.ToLower(lang[1:]))
		} else {
			templateName = "PasswordResetConfirmation"
		}
	}

	return
}

// main sets up the Lambda function entry point by registering the handleRequest
// function as the handler for Lambda events. The handleRequest function is designed
// to process SQSEvent messages and use AWS SES for sending confirmation emails
// to users based on the contents of these messages.
//
// The lambda.Start function call is responsible for bootstrapping the Lambda runtime
// and ensures that incoming AWS Lambda events are directed to the handleRequest function.
func main() {
	lambda.Start(handleRequest)
}
