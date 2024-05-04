package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/ctrlzdtoofar/verdure-lambdas/internal/mdl"
	"github.com/ctrlzdtoofar/verdure-lambdas/internal/settings"
	"testing"
)

// Add the test ses template:
// aws ses create-template --cli-input-json file://templates/test.json
// aws ses update-template --cli-input-json file://templates/test.json
// aws ses delete-template --template-name TestMessage

func TestNoParamTemplateEmailSend(t *testing.T) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-2"),
	)
	if err != nil {
		t.Errorf("failed to load default config, %v", err)
	}

	sesSvc := sesv2.NewFromConfig(cfg)
	templateName := "TestMessage"
	templateData := fmt.Sprintf("{\"url\": \"%s\"}", "testing123")

	toAddr := settings.GetEnv("TEST_EMAIL", "")
	if len(toAddr) == 0 {
		t.Errorf("Environment variable TEST_EMAIL is required")
	}

	noReplyAddr := settings.GetEnv("NOREPLY_EMAIL", "")
	if len(noReplyAddr) == 0 {
		noReplyAddr = toAddr
	}

	// Define Email Parameters
	input := &sesv2.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				toAddr,
			},
		},
		FromEmailAddress: &noReplyAddr,
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
	_, err = sesSvc.SendEmail(context.TODO(), input)
	if err != nil {
		t.Errorf("failed to send %s email, %v", templateName, err)
	}
}

func TestTemplatesAndEmailSend(t *testing.T) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-2"),
	)
	if err != nil {
		t.Errorf("failed to load default config, %v", err)
	}

	sesSvc := sesv2.NewFromConfig(cfg)

	toAddr := settings.GetEnv("TEST_EMAIL", "")
	if len(toAddr) == 0 {
		t.Errorf("Environment variable TEST_EMAIL is required")
	}

	noReplyAddr := settings.GetEnv("NOREPLY_EMAIL", "")
	if len(noReplyAddr) == 0 {
		noReplyAddr = toAddr
	}

	confirm := &mdl.UserConfirmation{
		ConfirmationType: mdl.NewUser,
		BaseUrl:          "/basetest",
		UserLoginID:      123,
		Email:            toAddr,
		Lang:             "en",
		Token:            "testing123",
		ExpiresAtMillis:  1700000000000,
	}

	err = sendEmail(context.TODO(), *sesSvc, *confirm, noReplyAddr)
	if len(toAddr) == 0 {
		t.Errorf("Failed to send email using template %v", err)
	}

	confirm.Lang = "es"
	err = sendEmail(context.TODO(), *sesSvc, *confirm, noReplyAddr)
	if len(toAddr) == 0 {
		t.Errorf("Failed to send email using template %v", err)
	}

	confirm.Lang = "de"
	err = sendEmail(context.TODO(), *sesSvc, *confirm, noReplyAddr)
	if len(toAddr) == 0 {
		t.Errorf("Failed to send email using template %v", err)
	}
}

func TestHandleRequest(t *testing.T) {

	toAddr := settings.GetEnv("TEST_EMAIL", "")
	if len(toAddr) == 0 {
		t.Errorf("Environment variable TEST_EMAIL is required")
	}

	confirm := &mdl.UserConfirmation{
		ConfirmationType: mdl.NewUser,
		BaseUrl:          "http://localhost:3000",
		UserLoginID:      123,
		Email:            toAddr,
		Lang:             "es",
		Token:            "testing123",
		ExpiresAtMillis:  19000000000,
	}

	// Serialize UserConfirmation to JSON
	confirmationJSON, err := json.Marshal(confirm)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Create an SQS event
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(confirmationJSON),
			},
		},
	}

	err = handleRequest(context.TODO(), sqsEvent)
	if err != nil {
		t.Errorf("Failed to invoke handler with test scenario, data %v, error %v", confirm, err)
	}

	confirm.ConfirmationType = mdl.ResetPassword
	// Serialize UserConfirmation to JSON
	confirmationJSON, err = json.Marshal(confirm)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Create an SQS event
	sqsEvent = events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(confirmationJSON),
			},
		},
	}

	err = handleRequest(context.TODO(), sqsEvent)
	if err != nil {
		t.Errorf("Failed to invoke handler with test scenario, data %v, error %v", confirm, err)
	}

}
