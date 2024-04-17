package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/heather92115/verdure-lambdas/internal/settings"
	"testing"
)

// Add the test ses template:
// aws ses create-template --cli-input-json file://templates/test.json
// aws ses update-template --cli-input-json file://templates/test.json
// aws ses delete-template --template-name TestMessage
//

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
