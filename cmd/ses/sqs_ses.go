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
				Name:  aws.String("env"),
				Value: aws.String("local"),
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
				strings.ToUpper(lang)+strings.ToLower(lang))
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
